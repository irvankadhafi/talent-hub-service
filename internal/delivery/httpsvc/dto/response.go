package dto

// LoginResponse for login response data.
type LoginResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
}

// Response is a generic structure for standard API responses.
type Response[T any] struct {
	Data    T      `json:"data,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// NewSuccessResponse creates a new success response.
func NewSuccessResponse[T any](data T, message string) Response[T] {
	return Response[T]{
		Data:    data,
		Message: message,
		Success: true,
	}
}

// NewErrorResponse creates a new error response.
func NewErrorResponse[T any](message string) Response[T] {
	return Response[T]{
		Message: message,
		Success: false,
	}
}
