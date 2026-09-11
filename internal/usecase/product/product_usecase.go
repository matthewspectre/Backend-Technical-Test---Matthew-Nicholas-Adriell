package product

import (
	"context"
	"errors"
	"strings"

	entity "be_evindo/internal/entity/product"
	repo "be_evindo/internal/repository/product"
)

type Usecase interface {
	Create(ctx context.Context, data *entity.Product) error
	FindByID(ctx context.Context, id int) (*entity.Product, error)
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
