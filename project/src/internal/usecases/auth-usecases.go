package usecases

import (
	"context"

	"gin-alpine/src/internal/domain/auth"
	"gin-alpine/src/internal/domain/users"
	"gin-alpine/src/pkg/utils"

	"gin-alpine/src/internal/sqlc/gen"
)

type AuthUsecases struct {
	usersRepo users.Repository
}

func NewAuthUsecases(usersRepo users.Repository) *AuthUsecases {
	return &AuthUsecases{usersRepo}
}

func (u *AuthUsecases) Login(ctx context.Context, input auth.LoginInput) (*gen.FindUserByEmailRow, error) {
	user, err := u.usersRepo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, users.ErrUserNotFound
	}

	err = utils.ValidatePassword(input.Password, user.Password)
	if err != nil {
		return nil, auth.ErrInvalidCredentials
	}
	return user, nil
}
