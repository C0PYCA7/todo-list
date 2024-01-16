package postgres

import "errors"

var (
	ErrUserNotFound    = errors.New("signIn not found")
	ErrUserNameExists  = errors.New("signIn with this username already exists")
	ErrUserLoginExists = errors.New("signIn with this login already exists")
)
