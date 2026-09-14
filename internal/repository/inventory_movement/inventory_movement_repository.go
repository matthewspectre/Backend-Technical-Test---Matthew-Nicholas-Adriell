package inventory_movement

import (
	"context"

	inventorymovemententity "be_evindo/internal/entity/inventory_movement"
)

type InventoryMovementRepository interface {
	FindAll(ctx context.Context, productID *int, warehouseID *int, reference string) ([]*inventorymovemententity.InventoryMovement, error)
}
