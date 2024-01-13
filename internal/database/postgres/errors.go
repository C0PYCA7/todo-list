package postgres

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserNameExists  = errors.New("user with this username already exists")
	ErrUserLoginExists = errors.New("user with this login already exists")
)
