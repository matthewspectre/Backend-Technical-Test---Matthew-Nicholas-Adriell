package usecase

import (
	"context"
	"errors"
	"strings"

	"be_evindo/internal/entity/user"
)

type UserUsecase struct {
	repository user.Repository
}

func NewUserUsecase(repository user.Repository) *UserUsecase {
	return &UserUsecase{repository: repository}
}

func (usecase *UserUsecase) FindByID(ctx context.Context, id int64) (*user.User, error) {
	if id <= 0 {
		return nil, errors.New("user id must be greater than zero")
	}
	return usecase.repository.FindByID(ctx, id)
}

func (usecase *UserUsecase) Create(ctx context.Context, entity *user.User) error {
	entity.Name = strings.TrimSpace(entity.Name)
	entity.Email = strings.TrimSpace(entity.Email)
	if entity.Name == "" || entity.Email == "" || entity.Password == "" {
		return errors.New("name, email, and password are required")
	}
	return usecase.repository.Create(ctx, entity)
}
