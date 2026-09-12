package product

import (
	"context"
	"errors"
	"strings"
	"time"

	entity "be_evindo/internal/entity/product"
	repo "be_evindo/internal/repository/product"
)

type Usecase interface {
	Create(ctx context.Context, data *entity.Product) error
	FindAll(ctx context.Context) ([]*entity.Product, error)
	FindByID(ctx context.Context, id int) (*entity.Product, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}

type usecase struct {
	repo repo.ProductRepository
}

func NewUsecase(repository repo.ProductRepository) Usecase {
	return &usecase{repo: repository}
}

func (usecase *usecase) Create(ctx context.Context, data *entity.Product) error {
	if data == nil {
		return errors.New("product is required")
	}
	data.SKU = strings.TrimSpace(data.SKU)
	data.Name = strings.TrimSpace(data.Name)
	data.Unit = strings.TrimSpace(data.Unit)
	if data.SKU == "" || data.Name == "" || data.Unit == "" {
		return errors.New("sku, name, and unit are required")
	}
	return usecase.repo.Create(ctx, data)
}

func (usecase *usecase) FindByID(ctx context.Context, id int) (*entity.Product, error) {
	if id <= 0 {
		return nil, errors.New("product id must be greater than zero")
	}
	return usecase.repo.FindByID(ctx, id)
}

func (usecase *usecase) FindAll(ctx context.Context) ([]*entity.Product, error) {
	return usecase.repo.FindAll(ctx)
}

func (usecase *usecase) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	if id <= 0 {
		return errors.New("product id must be greater than zero")
	}
	if len(updates) == 0 {
		return errors.New("at least one field is required")
	}
	updates["updated_at"] = time.Now().Format("2006-01-02 15:04:05.999999-07")
	if value, exists := updates["sku"]; exists {
		updates["sku"] = strings.TrimSpace(value.(string))
		if updates["sku"] == "" {
			return errors.New("sku must not be empty")
		}
	}
	if value, exists := updates["name"]; exists {
		updates["name"] = strings.TrimSpace(value.(string))
		if updates["name"] == "" {
			return errors.New("name must not be empty")
		}
	}
	if value, exists := updates["unit"]; exists {
		updates["unit"] = strings.TrimSpace(value.(string))
		if updates["unit"] == "" {
			return errors.New("unit must not be empty")
		}
	}
	return usecase.repo.Update(ctx, id, updates)
}

func (usecase *usecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("product id must be greater than zero")
	}
	return usecase.repo.Delete(ctx, id)
}
