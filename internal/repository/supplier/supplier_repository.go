package supplier

import (
	"context"

	supplierentity "be_evindo/internal/entity/supplier"
)

type SupplierRepository interface {
	FindAll(ctx context.Context) ([]*supplierentity.Supplier, error)
	FindByID(ctx context.Context, id int) (*supplierentity.Supplier, error)
	Create(ctx context.Context, entity *supplierentity.Supplier) error
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}
