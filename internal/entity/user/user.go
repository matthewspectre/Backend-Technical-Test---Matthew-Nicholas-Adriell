package user

import "context"

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt string
	UpdatedAt string
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}
