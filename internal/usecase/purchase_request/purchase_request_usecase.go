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
	Update(ctx context.Context, id int64, warehouseID *int, status *string, items []*entity.PurchaseRequestItem) error
	Approve(ctx context.Context, id int64) error
	Reject(ctx context.Context, id int64) error
	FindAll(ctx context.Context, status string) ([]*entity.PurchaseRequest, error)
	FindByID(ctx context.Context, id int64) (*entity.PurchaseRequest, error)
}

var ErrItemsRequired = errors.New("at least one item is required")
var ErrInvalidItemQuantity = errors.New("item quantity must be greater than zero")
var ErrInvalidStatus = errors.New("status must be DRAFT, SUBMITTED, APPROVED, or REJECTED")
var ErrStatusDataNotFound = errors.New("no purchase requests found for the selected status")
var ErrUpdateFieldsRequired = errors.New("at least one field is required")
var ErrNotDraftNotEditable = errors.New("only draft purchase requests can be edited")
var ErrApprovalRequiresSubmitted = errors.New("only submitted purchase requests can be approved")
var ErrRejectionRequiresSubmitted = errors.New("only submitted purchase requests can be rejected")

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

func (usecase *usecase) FindAll(ctx context.Context, status string) ([]*entity.PurchaseRequest, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != "" && status != "DRAFT" && status != "SUBMITTED" && status != "APPROVED" && status != "REJECTED" {
		return nil, ErrInvalidStatus
	}
	data, err := usecase.repository.FindAll(ctx, status)
	if err != nil {
		return nil, err
	}
	if status != "" && len(data) == 0 {
		return nil, ErrStatusDataNotFound
	}
	return data, nil
}

func (usecase *usecase) Update(ctx context.Context, id int64, warehouseID *int, status *string, items []*entity.PurchaseRequestItem) error {
	if id <= 0 {
		return errors.New("purchase request id must be greater than zero")
	}
	existing, err := usecase.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.Status != "DRAFT" {
		return ErrNotDraftNotEditable
	}
	if warehouseID == nil && status == nil && items == nil {
		return ErrUpdateFieldsRequired
	}
	if warehouseID != nil {
		if *warehouseID <= 0 {
			return errors.New("warehouse_id must be greater than zero")
		}
		warehouse, err := usecase.warehouseRepository.FindByID(ctx, *warehouseID)
		if err != nil {
			return err
		}
		if warehouse == nil || warehouse.IsActive != 1 {
			return errors.New("inactive warehouses cannot be used for purchase requests")
		}
	}
	if status != nil {
		*status = strings.ToUpper(strings.TrimSpace(*status))
		if *status != "DRAFT" && *status != "SUBMITTED" && *status != "APPROVED" && *status != "REJECTED" {
			return ErrInvalidStatus
		}
	}
	if items != nil {
		if len(items) == 0 {
			return ErrItemsRequired
		}
		seenProducts := make(map[int]struct{}, len(items))
		for _, item := range items {
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
	}
	return usecase.repository.Update(ctx, id, warehouseID, status, items)
}

func (usecase *usecase) FindByID(ctx context.Context, id int64) (*entity.PurchaseRequest, error) {
	if id <= 0 {
		return nil, errors.New("purchase request id must be greater than zero")
	}
	return usecase.repository.FindByID(ctx, id)
}

func (usecase *usecase) Approve(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("purchase request id must be greater than zero")
	}
	request, err := usecase.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if request.Status != "SUBMITTED" {
		return ErrApprovalRequiresSubmitted
	}
	return usecase.repository.Approve(ctx, id)
}

func (usecase *usecase) Reject(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("purchase request id must be greater than zero")
	}
	request, err := usecase.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if request.Status != "SUBMITTED" {
		return ErrRejectionRequiresSubmitted
	}
	return usecase.repository.Reject(ctx, id)
}
