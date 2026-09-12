package user

import (
	"context"

	entity "be_evindo/internal/entity/user"
	model "be_evindo/internal/model/user"
	repo "be_evindo/internal/repository/user"

	"gorm.io/gorm"
)

type RepositoryPostgre struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.UserRepository {
	return &RepositoryPostgre{db: db}
}

func toEntity(userModel *model.UserModel) *entity.User {
	if userModel == nil {
		return nil
	}
	return &entity.User{
		ID:        userModel.ID,
		Name:      userModel.Name,
		Email:     userModel.Email,
		Password:  userModel.Password,
		Role:      userModel.Role,
		CreatedAt: userModel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: userModel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, entityUser *entity.User) error {
	userModel := &model.UserModel{
		Name:     entityUser.Name,
		Email:    entityUser.Email,
		Password: entityUser.Password,
		Role:     entityUser.Role,
	}
	if err := repository.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return err
	}
	entityUser.ID = userModel.ID
	return nil
}

func (repository *RepositoryPostgre) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	userModel := &model.UserModel{}
	if err := repository.db.WithContext(ctx).Where("email = ?", email).First(userModel).Error; err != nil {
		return nil, err
	}
	return toEntity(userModel), nil
}
