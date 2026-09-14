package goods_receipt

type GoodsReceiptCreateRequest struct {
	Items []GoodsReceiptItemCreateRequest `json:"items"`
}

type GoodsReceiptItemCreateRequest struct {
	ProductID        int `json:"product_id"`
	ReceivedQuantity int `json:"received_quantity"`
}

type GoodsReceiptResponse struct {
	ID              int64                      `json:"id"`
	ReceiptNumber   string                     `json:"receipt_number"`
	PurchaseOrderID int64                      `json:"purchase_order_id"`
	WarehouseID     int                        `json:"warehouse_id"`
	ReceivedBy      int64                      `json:"received_by"`
	Status          string                     `json:"status"`
	Items           []GoodsReceiptItemResponse `json:"items"`
}

type GoodsReceiptItemResponse struct {
	ID                  int64 `json:"id"`
	PurchaseOrderItemID int64 `json:"purchase_order_item_id"`
	ProductID           int   `json:"product_id"`
	ReceivedQuantity    int   `json:"received_quantity"`
}

type GoodsReceiptListResponse struct {
	Data      []GoodsReceiptResponse `json:"data"`
	TotalData int                    `json:"total_data"`
}
