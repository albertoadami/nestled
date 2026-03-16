package errors

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrPasswordTooWeak       = errors.New("password is too weak")
	CredentialsInvalid       = errors.New("invalid credentials")
	ErrInvalidToken          = errors.New("invalid token")
	ErrGeneratingToken       = errors.New("error generating token")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidTokenType      = errors.New("invalid token provided")
	ErrNoRowsAffected        = errors.New("no rows affected by operation")
)
