package purchase_order

type PurchaseOrderCreateRequest struct {
	SupplierID int `json:"supplier_id"`
}

type PurchaseOrderResponse struct {
	ID                int64  `json:"id"`
	PONumber          string `json:"po_number"`
	PurchaseRequestID int64  `json:"purchase_request_id"`
	SupplierID        int    `json:"supplier_id"`
	SupplierName      string `json:"supplier_name"`
	WarehouseID       int    `json:"warehouse_id"`
	WarehouseName     string `json:"warehouse_name"`
	Status            string `json:"status"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}
