// Package auth
package auth

import "gin-alpine/src/pkg/utils"

type LoginInputDTO struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=4,max=100"`
}

func GetLoginCustomMessages(t *utils.Translator) *map[string]string {
	return &map[string]string{
		"Email.required":    t.T("errors.user.email_required", nil),
		"Email.email":       t.T("errors.user.valid_email", nil),
		"Password.required": t.T("errors.user.password_required", nil),
		"Password.min":      t.T("errors.user.password_min_length", nil),
	}
}
