package inventory_movement

import (
	"context"

	entity "be_evindo/internal/entity/inventory_movement"
	repo "be_evindo/internal/repository/inventory_movement"
)

type Usecase interface {
	FindAll(ctx context.Context, productID *int, warehouseID *int, reference string) ([]*entity.InventoryMovement, error)
}

type usecase struct {
	repository repo.InventoryMovementRepository
}

func NewUsecase(repository repo.InventoryMovementRepository) Usecase {
	return &usecase{repository: repository}
}

func (usecase *usecase) FindAll(ctx context.Context, productID *int, warehouseID *int, reference string) ([]*entity.InventoryMovement, error) {
	return usecase.repository.FindAll(ctx, productID, warehouseID, reference)
}
