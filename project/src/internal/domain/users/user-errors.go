package users

import (
	"errors"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrRoleNotAllowed     = errors.New("no permission to perform this action with your role")
	ErrInvalidRole        = errors.New("not a valid role to set")
)
