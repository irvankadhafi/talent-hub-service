package usecase

import (
	"context"
	"errors"
	"github.com/irvankadhafi/talent-hub-service/auth"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CandidateAutherAdapter adapter for auth.CandidateAuthenticator
type CandidateAutherAdapter struct {
	authUsecase model.AuthUsecase
}

// NewCandidateAutherAdapter constructor
func NewCandidateAutherAdapter(authUsecase model.AuthUsecase) *CandidateAutherAdapter {
	return &CandidateAutherAdapter{
		authUsecase: authUsecase,
	}
}

// AuthenticateToken authenticate access token
func (a *CandidateAutherAdapter) AuthenticateToken(ctx context.Context, accessToken string) (*auth.Candidate, error) {
	user, err := a.authUsecase.AuthenticateToken(ctx, accessToken)
	if errors.Is(err, ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, ErrAccessTokenExpired) {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if err != nil {
		return nil, err
	}

	return newAuthCandidate(user), nil
}

func newAuthCandidate(candidate *model.Candidate) *auth.Candidate {
	if candidate == nil {
		return nil
	}
	return &auth.Candidate{
		ID:        candidate.ID,
		SessionID: candidate.SessionID,
	}
}
