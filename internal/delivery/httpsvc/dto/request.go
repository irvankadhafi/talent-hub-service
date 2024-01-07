package dto

// LoginRequest for login request data.
type LoginRequest struct {
	Identifier string `json:"identifier"` // Can be either email or phone
	Password   string `json:"password"`
}
