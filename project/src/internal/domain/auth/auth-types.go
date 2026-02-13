// Package auth
package auth

import (
	"context"
	"fmt"
	"time"

	"gin-alpine/src/internal/infra/redis"

	"github.com/go-redis/cache/v9"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token        string
	RefreshToken string
}

type LogoutInput struct {
	AccessToken string
}

type GetRefreshTokenInput struct {
	RefreshToken string
}

type Role int

const (
	RoleCustomer Role = 1
	RoleManager  Role = 2
	RoleAdmin    Role = 3
	RoleDev      Role = 4
)

func (r Role) String() string {
	switch r {
	case RoleCustomer:
		return "CUSTOMER"
	case RoleManager:
		return "MANAGER"
	case RoleAdmin:
		return "ADMIN"
	case RoleDev:
		return "DEV"
	default:
		return "UNKNOWN"
	}
}

type UserAuth struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
}

func StoreUserAuth(
	ctx context.Context,
	r *redis.RedisClient,
	user UserAuth,
	ttl time.Duration,
) error {
	key := fmt.Sprintf("auth:user:%d", user.ID)

	return r.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: &user,
		TTL:   ttl,
	})
}
