package inventory

import (
	"context"

	inventoryentity "be_evindo/internal/entity/inventory"
)

type InventoryRepository interface {
	FindAll(ctx context.Context) ([]*inventoryentity.Inventory, error)
	FindByID(ctx context.Context, id int) (*inventoryentity.Inventory, error)
	FindByProduct(ctx context.Context, productID int) ([]*inventoryentity.Inventory, error)
	FindByWarehouse(ctx context.Context, warehouseID int) ([]*inventoryentity.Inventory, error)
	Create(ctx context.Context, entity *inventoryentity.Inventory) error
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}
