package errors

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrPasswordTooWeak       = errors.New("password is too weak")
)
