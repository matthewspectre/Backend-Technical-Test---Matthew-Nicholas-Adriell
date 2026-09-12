package warehouse

import (
	"context"
	"errors"
	"strings"
	"time"

	entity "be_evindo/internal/entity/warehouse"
	repo "be_evindo/internal/repository/warehouse"
)

type Usecase interface {
	Create(ctx context.Context, data *entity.Warehouse) error
	FindAll(ctx context.Context) ([]*entity.Warehouse, error)
	FindByID(ctx context.Context, id int) (*entity.Warehouse, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}

type usecase struct {
	repo repo.WarehouseRepository
}

func NewUsecase(repository repo.WarehouseRepository) Usecase {
	return &usecase{repo: repository}
}

func (usecase *usecase) Create(ctx context.Context, data *entity.Warehouse) error {
	if data == nil {
		return errors.New("warehouse is required")
	}
	data.Code = strings.ToUpper(strings.TrimSpace(data.Code))
	data.Name = strings.TrimSpace(data.Name)
	data.Location = strings.TrimSpace(data.Location)
	if data.Code == "" || data.Name == "" || data.Location == "" {
		return errors.New("code, name, and location are required")
	}
	return usecase.repo.Create(ctx, data)
}

func (usecase *usecase) FindAll(ctx context.Context) ([]*entity.Warehouse, error) {
	return usecase.repo.FindAll(ctx)
}

func (usecase *usecase) FindByID(ctx context.Context, id int) (*entity.Warehouse, error) {
	if id <= 0 {
		return nil, errors.New("warehouse id must be greater than zero")
	}
	return usecase.repo.FindByID(ctx, id)
}

func (usecase *usecase) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	if id <= 0 {
		return errors.New("warehouse id must be greater than zero")
	}
	if len(updates) == 0 {
		return errors.New("at least one field is required")
	}
	updates["updated_at"] = time.Now().Format("2006-01-02 15:04:05.999999-07")
	for _, field := range []string{"code", "name", "location"} {
		if value, exists := updates[field]; exists {
			text, ok := value.(string)
			if !ok || strings.TrimSpace(text) == "" {
				return errors.New(field + " must not be empty")
			}
			updates[field] = strings.TrimSpace(text)
		}
	}
	if value, exists := updates["code"]; exists {
		updates["code"] = strings.ToUpper(value.(string))
	}
	return usecase.repo.Update(ctx, id, updates)
}

func (usecase *usecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("warehouse id must be greater than zero")
	}
	return usecase.repo.Delete(ctx, id)
}
