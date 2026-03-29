package services

import "errors"

var (
	ErrInvalidAuthInput       = errors.New("invalid auth input")
	ErrAuthUserAlreadyExists  = errors.New("auth user already exists")
	ErrAuthInvalidCredentials = errors.New("invalid auth credentials")
	ErrAuthUserDisabled       = errors.New("auth user disabled")
)
