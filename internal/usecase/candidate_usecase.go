package usecase

import (
	"context"
	"github.com/irvankadhafi/talent-hub-service/internal/helper"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
)

type candidateUsecase struct {
	candidateRepo model.CandidateRepository
}

func NewCandidateUsecase(candidateRepo model.CandidateRepository) model.CandidateUsecase {
	return &candidateUsecase{
		candidateRepo: candidateRepo,
	}
}

func (c *candidateUsecase) Create(ctx context.Context, input model.CreateCandidateInput) (*model.Candidate, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"fullName": input.FullName,
		"email":    input.Email,
		"phone":    input.Phone,
	})

	input.Email = helper.FormatEmail(input.Email)
	if err := input.ValidateAndFormat(); err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := c.checkCandidateExistence(ctx, input.Email, input.Phone); err != nil {
		logger.Error(err)
		return nil, err
	}

	cipherPwd, err := helper.HashString(input.Password)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	candidateInput := &model.Candidate{
		ID:       utils.GenerateID(),
		FullName: input.FullName,
		Email:    null.StringFrom(input.Email),
		Phone:    null.StringFrom(input.Phone),
		Gender:   input.Gender,
		Password: cipherPwd,
	}

	if err := c.candidateRepo.Create(ctx, candidateInput); err != nil {
		logger.Error(err)
		return nil, err
	}

	return c.FindByID(ctx, candidateInput.ID)
}

func (c *candidateUsecase) FindByID(ctx context.Context, id int64) (*model.Candidate, error) {
	candidate, err := c.candidateRepo.FindByID(ctx, id)
	if err != nil {
		logrus.WithField("id", id).Error(err)
		return nil, err
	}

	if candidate == nil {
		return nil, ErrNotFound
	}

	return candidate, nil
}

func (c *candidateUsecase) checkCandidateExistence(ctx context.Context, email, phone string) error {
	if email != "" {
		if _, err := c.findCandidate(ctx, "email", email); err != nil {
			return err
		}
	}

	if phone != "" {
		if _, err := c.findCandidate(ctx, "phone", phone); err != nil {
			return err
		}
	}

	return nil
}

func (c *candidateUsecase) findCandidate(ctx context.Context, field, value string) (*model.Candidate, error) {
	var candidate *model.Candidate
	var err error

	switch field {
	case "email":
		candidate, err = c.candidateRepo.FindUnscopedByEmail(ctx, value)
	case "phone":
		candidate, err = c.candidateRepo.FindUnscopedByPhone(ctx, value)
	}

	switch {
	case candidate != nil:
		return nil, ErrDuplicateCandidate
	case err == ErrNotFound:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return candidate, nil
}
