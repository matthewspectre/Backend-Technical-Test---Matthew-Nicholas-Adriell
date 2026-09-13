package purchase_order

import (
	"context"
	"errors"

	entity "be_evindo/internal/entity/purchase_order"
	repo "be_evindo/internal/repository/purchase_order"
	purchaserequestrepo "be_evindo/internal/repository/purchase_request"
	supplierrepo "be_evindo/internal/repository/supplier"
	"gorm.io/gorm"
)

type Usecase interface {
	CreateFromPurchaseRequest(ctx context.Context, purchaseRequestID int64, supplierID int) (*entity.PurchaseOrder, error)
	FindByID(ctx context.Context, id int64) (*entity.PurchaseOrder, error)
}

var ErrPurchaseRequestNotApproved = errors.New("only approved purchase requests can create purchase orders")
var ErrPurchaseOrderAlreadyExists = errors.New("purchase order already exists for this purchase request")

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
		return nil, errors.New("inactive suppliers cannot be used for purchase orders")
	}
	data := &entity.PurchaseOrder{
		PurchaseRequestID: purchaseRequestID,
		SupplierID:        supplierID,
		WarehouseID:       request.WarehouseID,
		Status:            "DRAFT",
	}
	if err := usecase.repository.Create(ctx, data); err != nil {
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
