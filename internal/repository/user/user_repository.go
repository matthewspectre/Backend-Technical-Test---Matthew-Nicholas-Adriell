package user

import (
	"context"

	userentity "be_evindo/internal/entity/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *userentity.User) error
	FindByEmail(ctx context.Context, email string) (*userentity.User, error)
}
