package usecase

import (
	"context"
	"github.com/irvankadhafi/talent-hub-service/internal/config"
	"github.com/irvankadhafi/talent-hub-service/internal/helper"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"github.com/sirupsen/logrus"
	"time"
)

type authUsecase struct {
	candidateRepo model.CandidateRepository
	sessionRepo   model.SessionRepository
}

func NewAuthUsecase(
	candidateRepo model.CandidateRepository,
	sessionRepo model.SessionRepository,
) model.AuthUsecase {
	return &authUsecase{
		candidateRepo: candidateRepo,
		sessionRepo:   sessionRepo,
	}
}

// LoginByIdentifierPassword is refactored to handle both email and phone logins.
func (a *authUsecase) LoginByIdentifierPassword(ctx context.Context, req model.LoginRequest) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.DumpIncomingContext(ctx),
		"identifier": req.Identifier,
	})

	if err := req.Validate(); err != nil {
		logger.Error(err)
		return nil, err
	}

	var candidate *model.Candidate
	var err error
	if helper.ValidateEmail(req.Identifier) {
		candidate, err = a.findCandidateByEmail(ctx, req.Identifier)
	} else {
		if err := helper.RemoveLeadingZeroPhoneNumber(&req.Identifier); err != nil {
			logger.Error(err)
			return nil, err
		}

		if err := helper.FormatPhoneNumberWithCountryCode(&req.Identifier, "ID"); err != nil {
			logger.Error(err)
			return nil, err
		}

		candidate, err = a.findCandidateByPhone(ctx, req.Identifier)
	}
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return a.authenticateAndCreateSession(ctx, candidate, req)
}

func (a *authUsecase) AuthenticateToken(ctx context.Context, accessToken string) (*model.Candidate, error) {
	session, err := a.sessionRepo.FindByToken(ctx, model.AccessToken, accessToken)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if session == nil {
		return nil, ErrNotFound
	}

	if session.IsAccessTokenExpired() {
		return nil, ErrAccessTokenExpired
	}

	candidate, err := a.candidateRepo.FindByID(ctx, session.CandidateID)
	if err != nil {
		logrus.WithField("candidateID", session.CandidateID).Error(err)
		return nil, err
	}

	if candidate == nil {
		return nil, ErrNotFound
	}

	candidate.SessionID = session.ID

	return candidate, nil
}

// DeleteSessionByID deletes session by id.
func (a *authUsecase) DeleteSessionByID(ctx context.Context, sessionID int64) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"sessionID": utils.Dump(sessionID),
	})

	session, err := a.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if session == nil {
		return ErrNotFound
	}

	err = a.sessionRepo.Delete(ctx, session)
	if err != nil {
		logger.Error(err)
	}

	return err
}

// RefreshToken refresh the user's access and refresh token
func (a *authUsecase) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"ipAddress": req.IPAddress,
		"userAgent": req.UserAgent,
	})

	session, err := a.sessionRepo.FindByToken(ctx, model.RefreshToken, req.RefreshToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if session == nil {
		logger.Error(ErrNotFound)
		return nil, ErrNotFound
	}

	candidate, err := a.candidateRepo.FindByID(ctx, session.CandidateID)
	switch {
	case err != nil:
		logger.WithField("candidateID", session.CandidateID).Error(err)
		return nil, err
	case candidate == nil:
		logger.WithField("candidateID", session.CandidateID).Error(ErrNotFound)
		return nil, ErrNotFound
	}

	// old session is used to delete the old session cache
	oldSess := *session

	if session.RefreshTokenExpiredAt.Before(time.Now()) {
		logger.Error(ErrRefreshTokenExpired)
		return nil, ErrRefreshTokenExpired
	}

	newAccessToken, err := GenerateToken(a.sessionRepo, session.CandidateID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	newRefreshToken, err := GenerateToken(a.sessionRepo, session.CandidateID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	session.AccessToken = newAccessToken
	session.RefreshToken = newRefreshToken
	session.IPAddress = req.IPAddress
	session.UserAgent = req.UserAgent
	session.Latitude = req.Latitude
	session.Longitude = req.Longitude

	now := time.Now()
	session.AccessTokenExpiredAt = now.Add(config.AccessTokenDuration())
	session.RefreshTokenExpiredAt = now.Add(config.RefreshTokenDuration())

	session, err = a.sessionRepo.RefreshToken(ctx, &oldSess, session)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return session, nil
}

func (a *authUsecase) authenticateAndCreateSession(ctx context.Context, candidate *model.Candidate, req model.LoginRequest) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":         utils.DumpIncomingContext(ctx),
		"candidateID": candidate.ID,
		"userAgent":   req.UserAgent,
		"ipAddress":   req.IPAddress,
		"latitude":    req.Latitude,
		"longitude":   req.Longitude,
	})

	// Find the password associated with the candidate.
	cipherPass, err := a.candidateRepo.FindPasswordByID(ctx, candidate.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if cipherPass == nil {
		return nil, ErrUnauthorized
	}

	// Check if the provided password matches.
	if !helper.IsHashedStringMatch([]byte(req.PlainPassword), cipherPass) {
		return nil, ErrUnauthorized
	}

	// Generate access and refresh tokens.
	accessToken, err := GenerateToken(a.sessionRepo, candidate.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	refreshToken, err := GenerateToken(a.sessionRepo, candidate.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &model.Session{
		ID:                    utils.GenerateID(),
		CandidateID:           candidate.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  now.Add(config.AccessTokenDuration()),
		RefreshTokenExpiredAt: now.Add(config.RefreshTokenDuration()),
		IPAddress:             req.IPAddress,
		UserAgent:             req.UserAgent,
		Latitude:              req.Latitude,
		Longitude:             req.Longitude,
	}

	if err = a.sessionRepo.Create(ctx, session); err != nil {
		logger.Error(err)
		return nil, err
	}

	// TODO: implement delete session by worker

	return session, nil
}

func (a *authUsecase) findCandidateByEmail(ctx context.Context, email string) (*model.Candidate, error) {
	candidate, err := a.candidateRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if candidate == nil {
		return nil, ErrNotFound
	}

	return candidate, nil
}

func (a *authUsecase) findCandidateByPhone(ctx context.Context, phone string) (*model.Candidate, error) {
	candidate, err := a.candidateRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if candidate == nil {
		return nil, ErrNotFound
	}

	return candidate, nil
}
