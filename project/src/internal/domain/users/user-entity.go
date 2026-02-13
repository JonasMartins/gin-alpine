// Package users ...
package users

import (
	"gin-alpine/src/internal/domain/role"
	base "gin-alpine/src/pkg/models"
)

type User struct {
	Base     base.Base
	Name     string
	Email    string
	Password string
	Role     *role.Role
	RoleType string
	Enabled  bool
}
