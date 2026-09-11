package product

import (
	"context"

	productentity "be_evindo/internal/entity/product"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id int) (*productentity.Product, error)
	Create(ctx context.Context, entity *productentity.Product) error
}
