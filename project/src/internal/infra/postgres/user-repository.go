package postgres

import (
	"context"
	"time"

	"gin-alpine/src/internal/domain/role"
	"gin-alpine/src/internal/domain/users"
	"gin-alpine/src/internal/sqlc/gen"
	"gin-alpine/src/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	R *PgRepository
}

func NewUserRepository(p *PgRepository) *UserRepository {
	return &UserRepository{R: p}
}

func (r *UserRepository) FindUserByID(ctx context.Context, id int) (*gen.FindUserByIDRow, error) {
	return nil, nil
}
func (r *UserRepository) UpdateUser(ctx context.Context, input users.UpdateUserInput) error {
	return nil
}
func (r *UserRepository) UpdateUserAdmin(ctx context.Context, input users.UpdateUserAdminInput) error {
	return nil
}
func (r *UserRepository) GetUsersWithTotal(ctx context.Context, limit, offset int32) ([]*gen.GetUsersWithTotalRow, error) {
	return nil, nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*gen.FindUserByEmailRow, error) {
	q := gen.New(r.R.DB)
	u, err := q.FindUserByEmail(ctx, email)
	return &u, err
}

func (r *UserRepository) CreateUser(ctx context.Context, input *users.CreateUserInput) (*users.CreateUserResult, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, utils.ErrInternal
	}
	roleID := int32(1)
	if input.RoleID != nil && *input.RoleID > 0 {
		validRole := role.PositionType(uint8(*input.RoleID))
		if role.IsValidPositionType(validRole) {
			roleID = *input.RoleID
		}
	}
	q := gen.New(r.R.DB)
	row, err := q.CreateUser(ctx, gen.CreateUserParams{
		Uuid:      uuid.New(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashPass),
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		RoleID:    int32(roleID),
	})
	if err != nil {
		return nil, utils.ErrUnexpected
	}
	return &users.CreateUserResult{
		ID: int(row),
	}, nil

}
