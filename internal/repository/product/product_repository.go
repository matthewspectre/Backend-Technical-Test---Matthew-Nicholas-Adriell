package product

import (
	"context"

	productentity "be_evindo/internal/entity/product"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]*productentity.Product, error)
	FindByID(ctx context.Context, id int) (*productentity.Product, error)
	Create(ctx context.Context, entity *productentity.Product) error
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}
