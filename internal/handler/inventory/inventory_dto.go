package inventory

type InventoryCreateRequest struct {
	ProductID   int `json:"product_id"`
	WarehouseID int `json:"warehouse_id"`
	Stock       int `json:"stock"`
}

type InventoryUpdateRequest struct {
	ProductID   *int `json:"product_id,omitempty"`
	WarehouseID *int `json:"warehouse_id,omitempty"`
	Stock       *int `json:"stock,omitempty"`
}

type InventoryResponse struct {
	ID            int     `json:"id"`
	ProductID     int     `json:"product_id"`
	ProductName   string  `json:"product_name"`
	WarehouseID   int     `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	Stock         int     `json:"stock"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
}

type InventoryListResponse struct {
	Data      []InventoryResponse `json:"data"`
	TotalData int                 `json:"total_data"`
}
