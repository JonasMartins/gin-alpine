// Package http
package http

import (
	"errors"
	"net/http"

	"gin-alpine/src/internal/domain/auth"
	"gin-alpine/src/internal/domain/users"
	"gin-alpine/src/pkg/utils"
)

type HTTPError struct {
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
	Status  int    `json:"-"`
}

func MapAuthErrorToHTTP(err error, t *utils.Translator) (int, HTTPError) {
	msg := t.T("errors.unxpected", nil)
	switch {
	case errors.Is(err, auth.ErrInvalidToken):
		msg = t.T("errors.auth.invalid_token", nil)
		return http.StatusBadRequest, HTTPError{Message: msg, Field: "access_token"}
	case errors.Is(err, utils.ErrRedisKeyNotFound):
		msg = t.T("errors.unxpected", nil)
		return http.StatusInternalServerError, HTTPError{Message: msg, Field: "password"}
	case errors.Is(err, auth.ErrInvalidCredentials):
		msg = t.T("errors.user.invalid_credentials", nil)
		return http.StatusUnauthorized, HTTPError{Message: msg, Field: "password"}
	case errors.Is(err, users.ErrUserNotFound):
		msg = t.T("errors.user.not_found", nil)
		return http.StatusNotFound, HTTPError{Message: msg, Field: "email"}
	default:
		return http.StatusInternalServerError, HTTPError{Message: msg}
	}
}

func MapUserErrorToHTTP(err error, t *utils.Translator) (int, HTTPError) {
	msg := t.T("errors.unxpected", nil)
	switch {
	case errors.Is(err, users.ErrInvalidRole):
		msg = t.T("errors.user.invalid_role", nil)
		return http.StatusForbidden, HTTPError{Message: msg}
	case errors.Is(err, users.ErrRoleNotAllowed):
		msg = t.T("errors.user.role_forbidden", nil)
		return http.StatusForbidden, HTTPError{Message: msg}
	case errors.Is(err, users.ErrUserNotFound):
		msg = t.T("errors.user.not_found", nil)
		return http.StatusNotFound, HTTPError{Message: msg, Field: "email"}
	case errors.Is(err, users.ErrEmailAlreadyExists):
		msg = t.T("errors.user.email_already_exists", nil)
		return http.StatusConflict, HTTPError{Message: msg, Field: "email"}
	case errors.Is(err, utils.ErrInternal):
		msg = t.T("errors.internal", nil)
		return http.StatusInternalServerError, HTTPError{Message: msg}
	case errors.Is(err, utils.ErrDBConnectionFailed):
		msg = t.T("errors.db_conn", nil)
		return http.StatusServiceUnavailable, HTTPError{Message: msg}
	default:
		return http.StatusInternalServerError, HTTPError{Message: msg}
	}
}
