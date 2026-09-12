package warehouse

import (
	"context"

	warehouseentity "be_evindo/internal/entity/warehouse"
)

type WarehouseRepository interface {
	FindAll(ctx context.Context) ([]*warehouseentity.Warehouse, error)
	FindByID(ctx context.Context, id int) (*warehouseentity.Warehouse, error)
	Create(ctx context.Context, entity *warehouseentity.Warehouse) error
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}
