package users

import (
	"context"

	"gin-alpine/src/internal/sqlc/gen"
)

type Repository interface {
	CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserResult, error)
	FindUserByEmail(ctx context.Context, email string) (*gen.FindUserByEmailRow, error)
	FindUserByID(ctx context.Context, id int) (*gen.FindUserByIDRow, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) error
	UpdateUserAdmin(ctx context.Context, input UpdateUserAdminInput) error
	GetUsersWithTotal(ctx context.Context, limit, offset int32) ([]*gen.GetUsersWithTotalRow, error)
}
