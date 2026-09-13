package purchase_order

import "time"

type PurchaseOrderModel struct {
	ID                int64     `gorm:"primaryKey;autoIncrement;column:id"`
	PONumber          string    `gorm:"column:po_number;default:(-)"`
	PurchaseRequestID int64     `gorm:"column:purchase_request_id"`
	SupplierID        int       `gorm:"column:supplier_id"`
	SupplierName      string    `gorm:"column:supplier_name;->"`
	WarehouseID       int       `gorm:"column:warehouse_id"`
	WarehouseName     string    `gorm:"column:warehouse_name;->"`
	Status            string    `gorm:"column:status"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (PurchaseOrderModel) TableName() string {
	return "purchase_orders"
}
