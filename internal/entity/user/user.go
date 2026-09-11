package user

import "context"

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Repository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, user *User) error
}

type Usecase interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, user *User) error
}
