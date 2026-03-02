package dto

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func NewErrorResponse(message string, details string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Details: details,
	}
}

const (
	ErrUsernameAlreadyExists = "username already exists"
	ErrEmailAlreadyExists    = "email already exists"
	ErrInvalidToken          = "invalid token"
)
