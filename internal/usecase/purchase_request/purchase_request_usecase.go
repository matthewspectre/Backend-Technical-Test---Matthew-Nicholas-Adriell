package purchase_request

import (
	"context"
	"errors"
	"strings"

	entity "be_evindo/internal/entity/purchase_request"
	repo "be_evindo/internal/repository/purchase_request"
	warehouserepo "be_evindo/internal/repository/warehouse"
)

type Usecase interface {
	Create(ctx context.Context, data *entity.PurchaseRequest) error
	FindAll(ctx context.Context) ([]*entity.PurchaseRequest, error)
	FindByID(ctx context.Context, id int64) (*entity.PurchaseRequest, error)
}

var ErrItemsRequired = errors.New("at least one item is required")
var ErrInvalidItemQuantity = errors.New("item quantity must be greater than zero")

type usecase struct {
	repository          repo.PurchaseRequestRepository
	warehouseRepository warehouserepo.WarehouseRepository
}

func NewUsecase(repository repo.PurchaseRequestRepository, warehouseRepository warehouserepo.WarehouseRepository) Usecase {
	return &usecase{repository: repository, warehouseRepository: warehouseRepository}
}

func (usecase *usecase) Create(ctx context.Context, data *entity.PurchaseRequest) error {
	if data == nil {
		return errors.New("purchase request is required")
	}
	if data.RequestedBy <= 0 {
		return errors.New("requested_by must be greater than zero")
	}
	data.RequestNumber = strings.TrimSpace(data.RequestNumber)
	if data.RequestNumber == "" {
		return errors.New("request_number is required")
	}
	data.Status = "DRAFT"
	if data.WarehouseID <= 0 {
		return errors.New("warehouse_id must be greater than zero")
	}
	if len(data.Items) == 0 {
		return ErrItemsRequired
	}
	seenProducts := make(map[int]struct{}, len(data.Items))
	for _, item := range data.Items {
		if item == nil || item.ProductID <= 0 {
			return errors.New("each item must have a valid product_id")
		}
		if item.Quantity <= 0 {
			return ErrInvalidItemQuantity
		}
		if _, exists := seenProducts[item.ProductID]; exists {
			return errors.New("a product may only appear once in items")
		}
		seenProducts[item.ProductID] = struct{}{}
	}
	warehouse, err := usecase.warehouseRepository.FindByID(ctx, data.WarehouseID)
	if err != nil {
		return err
	}
	if warehouse == nil || warehouse.IsActive != 1 {
		return errors.New("inactive warehouses cannot be used for purchase requests")
	}
	return usecase.repository.Create(ctx, data)
}

func (usecase *usecase) FindAll(ctx context.Context) ([]*entity.PurchaseRequest, error) {
	return usecase.repository.FindAll(ctx)
}

func (usecase *usecase) FindByID(ctx context.Context, id int64) (*entity.PurchaseRequest, error) {
	if id <= 0 {
		return nil, errors.New("purchase request id must be greater than zero")
	}
	return usecase.repository.FindByID(ctx, id)
}
