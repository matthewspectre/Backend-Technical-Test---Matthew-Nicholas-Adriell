package inventory_movement

import (
	"context"

	entity "be_evindo/internal/entity/inventory_movement"
	model "be_evindo/internal/model/inventory_movement"
	repo "be_evindo/internal/repository/inventory_movement"

	"gorm.io/gorm"
)

type RepositoryPostgre struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) repo.InventoryMovementRepository { return &RepositoryPostgre{db: db} }

func toEntity(data *model.InventoryMovementModel) *entity.InventoryMovement {
	if data == nil {
		return nil
	}
	return &entity.InventoryMovement{
		ID:           data.ID,
		WarehouseID:  data.WarehouseID,
		ProductID:    data.ProductID,
		MovementType: data.MovementType,
		Quantity:     data.Quantity,
		Reference:    data.Reference,
		CreatedAt:    data.CreatedAt,
	}
}

func (repository *RepositoryPostgre) FindAll(ctx context.Context, productID *int, warehouseID *int, reference string) ([]*entity.InventoryMovement, error) {
	query := repository.db.WithContext(ctx).Order("id DESC")
	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}
	if warehouseID != nil {
		query = query.Where("warehouse_id = ?", *warehouseID)
	}
	if reference != "" {
		query = query.Where("reference ILIKE ?", "%"+reference+"%")
	}
	var movements []*model.InventoryMovementModel
	if err := query.Find(&movements).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.InventoryMovement, 0, len(movements))
	for _, movement := range movements {
		result = append(result, toEntity(movement))
	}
	return result, nil
}
