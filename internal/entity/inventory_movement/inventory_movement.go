package inventory_movement

import "time"

type InventoryMovement struct {
	ID           int64     `json:"id"`
	WarehouseID  int       `json:"warehouse_id"`
	ProductID    int       `json:"product_id"`
	MovementType string    `json:"movement_type"`
	Quantity     int       `json:"quantity"`
	Reference    string    `json:"reference"`
	CreatedAt    time.Time `json:"created_at"`
}
