package repository

import (
	"context"

	"be_evindo/internal/entity/user"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*user.User, error)
	Create(ctx context.Context, entity *user.User) error
}
