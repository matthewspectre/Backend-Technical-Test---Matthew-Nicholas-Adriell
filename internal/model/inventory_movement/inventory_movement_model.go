package inventory_movement

import "time"

type InventoryMovementModel struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id"`
	WarehouseID  int       `gorm:"column:warehouse_id"`
	ProductID    int       `gorm:"column:product_id"`
	MovementType string    `gorm:"column:movement_type"`
	Quantity     int       `gorm:"column:quantity"`
	Reference    string    `gorm:"column:reference"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (InventoryMovementModel) TableName() string {
	return "inventory_movements"
}
