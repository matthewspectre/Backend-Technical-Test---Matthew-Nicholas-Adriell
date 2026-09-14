package inventory_movement

type InventoryMovementResponse struct {
	ID           int64  `json:"id"`
	WarehouseID  int    `json:"warehouse_id"`
	ProductID    int    `json:"product_id"`
	MovementType string `json:"movement_type"`
	Quantity     int    `json:"quantity"`
	Reference    string `json:"reference"`
	CreatedAt    string `json:"created_at"`
}

type InventoryMovementListResponse struct {
	Data      []InventoryMovementResponse `json:"data"`
	TotalData int                         `json:"total_data"`
}
