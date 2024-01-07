package model

import (
	"context"
)

type IdentifierType int

const (
	IdentifierTypeUnknown IdentifierType = iota
	IdentifierTypeEmail
	IdentifierTypePhone
)

// LoginRequest request
type LoginRequest struct {
	Identifier     string `json:"identifier" validate:"required,identifier"`
	PlainPassword  string `json:"plain_password" validate:"required,min=5"`
	IdentifierType IdentifierType
	UserAgent      string `json:"user_agent"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	IPAddress      string `json:"ip_address"`
}

// Validate validates the login input body.
func (c *LoginRequest) Validate() error {
	return validate.Struct(c)
}

// RefreshTokenRequest request
type RefreshTokenRequest struct {
	RefreshToken string
	UserAgent    string
	Latitude     string
	Longitude    string
	IPAddress    string
}

// AuthUsecase usecases
type AuthUsecase interface {
	LoginByIdentifierPassword(ctx context.Context, req LoginRequest) (*Session, error)
	AuthenticateToken(ctx context.Context, accessToken string) (*Candidate, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*Session, error)
	DeleteSessionByID(ctx context.Context, sessionID int64) error
}
