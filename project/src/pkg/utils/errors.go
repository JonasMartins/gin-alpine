package utils

import (
	"errors"
)

var (
	ErrInternal           = errors.New("internal error")
	ErrDBConnectionFailed = errors.New("unable to connect to the database")
	ErrUnexpected         = errors.New("unexpected error occurred")
	ErrRedisKeyNotFound   = errors.New("redis key not found")
)
