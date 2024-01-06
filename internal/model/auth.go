package model

import "context"

// LoginRequest request
type LoginRequest struct {
	Email         string
	Phone         string
	PlainPassword string
}

// AuthUsecase usecases
type AuthUsecase interface {
	LoginByEmailPassword(ctx context.Context, req LoginRequest) (*Candidate, error)
	LoginByPhonePassword(ctx context.Context, req LoginRequest) (*Candidate, error)
	AuthenticateToken(ctx context.Context, accessToken string) (*Candidate, error)
	DeleteSessionByID(ctx context.Context, sessionID int64) error
}
