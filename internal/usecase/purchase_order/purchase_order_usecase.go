package purchase_order

import (
	"context"
	"errors"
	"strings"

	entity "be_evindo/internal/entity/purchase_order"
	repo "be_evindo/internal/repository/purchase_order"
	purchaserequestrepo "be_evindo/internal/repository/purchase_request"
	supplierrepo "be_evindo/internal/repository/supplier"

	"gorm.io/gorm"
)

type Usecase interface {
	CreateFromPurchaseRequest(ctx context.Context, purchaseRequestID int64, supplierID int) (*entity.PurchaseOrder, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	FindAll(ctx context.Context, status string) ([]*entity.PurchaseOrder, error)
	FindByID(ctx context.Context, id int64) (*entity.PurchaseOrder, error)
	FindByPurchaseRequestID(ctx context.Context, purchaseRequestID int64) (*entity.PurchaseOrder, error)
}

var ErrPurchaseRequestNotApproved = errors.New("only approved purchase requests can create purchase orders")
var ErrPurchaseOrderAlreadyExists = errors.New("purchase order already exists for this purchase request")
var ErrInvalidPurchaseOrderStatus = errors.New("status must be DRAFT, ORDERED, PARTIALLY_RECEIVED, RECEIVED, or CANCELLED")
var ErrInactiveSupplier = errors.New("inactive supplier cannot be used for purchase orders")

type usecase struct {
	repository                repo.PurchaseOrderRepository
	purchaseRequestRepository purchaserequestrepo.PurchaseRequestRepository
	supplierRepository        supplierrepo.SupplierRepository
}

func NewUsecase(repository repo.PurchaseOrderRepository, purchaseRequestRepository purchaserequestrepo.PurchaseRequestRepository, supplierRepository supplierrepo.SupplierRepository) Usecase {
	return &usecase{repository: repository, purchaseRequestRepository: purchaseRequestRepository, supplierRepository: supplierRepository}
}

func (usecase *usecase) CreateFromPurchaseRequest(ctx context.Context, purchaseRequestID int64, supplierID int) (*entity.PurchaseOrder, error) {
	if purchaseRequestID <= 0 {
		return nil, errors.New("purchase request id must be greater than zero")
	}
	if supplierID <= 0 {
		return nil, errors.New("supplier_id must be greater than zero")
	}
	request, err := usecase.purchaseRequestRepository.FindByID(ctx, purchaseRequestID)
	if err != nil {
		return nil, err
	}
	if request.Status != "APPROVED" {
		return nil, ErrPurchaseRequestNotApproved
	}
	existing, err := usecase.repository.FindByPurchaseRequestID(ctx, purchaseRequestID)
	if err == nil && existing != nil {
		return nil, ErrPurchaseOrderAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	supplier, err := usecase.supplierRepository.FindByID(ctx, supplierID)
	if err != nil {
		return nil, err
	}
	if supplier == nil || supplier.IsActive != 1 {
		return nil, ErrInactiveSupplier
	}
	data := &entity.PurchaseOrder{
		PurchaseRequestID: purchaseRequestID,
		SupplierID:        supplierID,
		WarehouseID:       request.WarehouseID,
		Status:            "DRAFT",
	}
	items := make([]*entity.PurchaseOrderItem, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, &entity.PurchaseOrderItem{
			ProductID: item.ProductID, OrderedQuantity: item.Quantity,
			ReceivedQuantity: 0,
		})
	}
	if err := usecase.repository.Create(ctx, data, items); err != nil {
		return nil, err
	}
	return usecase.repository.FindByID(ctx, data.ID)
}

func (usecase *usecase) FindByID(ctx context.Context, id int64) (*entity.PurchaseOrder, error) {
	if id <= 0 {
		return nil, errors.New("purchase order id must be greater than zero")
	}
	return usecase.repository.FindByID(ctx, id)
}

func (usecase *usecase) FindByPurchaseRequestID(ctx context.Context, purchaseRequestID int64) (*entity.PurchaseOrder, error) {
	if purchaseRequestID <= 0 {
		return nil, errors.New("purchase request id must be greater than zero")
	}
	return usecase.repository.FindByPurchaseRequestID(ctx, purchaseRequestID)
}

func (usecase *usecase) FindAll(ctx context.Context, status string) ([]*entity.PurchaseOrder, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != "" && status != "DRAFT" && status != "ORDERED" && status != "PARTIALLY_RECEIVED" && status != "RECEIVED" && status != "CANCELLED" {
		return nil, ErrInvalidPurchaseOrderStatus
	}
	return usecase.repository.FindAll(ctx, status)
}

func (usecase *usecase) UpdateStatus(ctx context.Context, id int64, status string) error {
	if id <= 0 {
		return errors.New("purchase order id must be greater than zero")
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != "DRAFT" && status != "ORDERED" && status != "PARTIALLY_RECEIVED" && status != "RECEIVED" && status != "CANCELLED" {
		return ErrInvalidPurchaseOrderStatus
	}
	return usecase.repository.UpdateStatus(ctx, id, status)
}
