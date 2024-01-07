package usecase

import "errors"

// errors ...
var (
	ErrNotFound                   = errors.New("not found")
	ErrUnauthorized               = errors.New("unauthorized")
	ErrAccessTokenExpired         = errors.New("access token expired")
	ErrRefreshTokenExpired        = errors.New("refresh token expired")
	ErrLoginByEmailPasswordLocked = errors.New("user is locked from logging in using email and password")
	ErrPermissionDenied           = errors.New("permission denied")
	ErrDuplicateCandidate         = errors.New("candidate already exist")
)
