package go_ecom

import "errors"

var (
	ErrPasswordMismatch = errors.New("passwords do not match")
	ErrPasswordTooShort = errors.New("password is too short")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrUserExists       = errors.New("user already exists")
)
