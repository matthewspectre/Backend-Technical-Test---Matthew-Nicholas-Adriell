package purchase_order

import (
	"context"
	"testing"

	purchaseorderentity "be_evindo/internal/entity/purchase_order"
	purchaserequestentity "be_evindo/internal/entity/purchase_request"
	supplierentity "be_evindo/internal/entity/supplier"
)

type fakePurchaseOrderRepository struct {
	created bool
	order   *purchaseorderentity.PurchaseOrder
}

func (f *fakePurchaseOrderRepository) Create(ctx context.Context, data *purchaseorderentity.PurchaseOrder, items []*purchaseorderentity.PurchaseOrderItem) error {
	f.created = true
	return nil
}

func (f *fakePurchaseOrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return nil
}

func (f *fakePurchaseOrderRepository) FindAll(ctx context.Context, status string) ([]*purchaseorderentity.PurchaseOrder, error) {
	return nil, nil
}

func (f *fakePurchaseOrderRepository) FindByID(ctx context.Context, id int64) (*purchaseorderentity.PurchaseOrder, error) {
	return &purchaseorderentity.PurchaseOrder{ID: id, Status: "DRAFT"}, nil
}

func (f *fakePurchaseOrderRepository) FindByPurchaseRequestID(ctx context.Context, purchaseRequestID int64) (*purchaseorderentity.PurchaseOrder, error) {
	if f.order == nil {
		return nil, nil
	}
	return f.order, nil
}

type fakePurchaseRequestRepository struct {
	request *purchaserequestentity.PurchaseRequest
}

func (f *fakePurchaseRequestRepository) Create(ctx context.Context, data *purchaserequestentity.PurchaseRequest) error {
	return nil
}

func (f *fakePurchaseRequestRepository) Update(ctx context.Context, id int64, warehouseID *int, status *string, items []*purchaserequestentity.PurchaseRequestItem) error {
	return nil
}

func (f *fakePurchaseRequestRepository) Approve(ctx context.Context, id int64) error {
	return nil
}

func (f *fakePurchaseRequestRepository) Reject(ctx context.Context, id int64) error {
	return nil
}

func (f *fakePurchaseRequestRepository) FindAll(ctx context.Context, status string) ([]*purchaserequestentity.PurchaseRequest, error) {
	return nil, nil
}

func (f *fakePurchaseRequestRepository) FindByID(ctx context.Context, id int64) (*purchaserequestentity.PurchaseRequest, error) {
	if f.request == nil {
		return nil, nil
	}
	return f.request, nil
}

type fakeSupplierRepository struct {
	supplier *supplierentity.Supplier
}

func (f *fakeSupplierRepository) FindAll(ctx context.Context) ([]*supplierentity.Supplier, error) {
	return nil, nil
}

func (f *fakeSupplierRepository) FindByID(ctx context.Context, id int) (*supplierentity.Supplier, error) {
	if f.supplier == nil {
		return nil, nil
	}
	return f.supplier, nil
}

func (f *fakeSupplierRepository) Create(ctx context.Context, entity *supplierentity.Supplier) error {
	return nil
}

func (f *fakeSupplierRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	return nil
}

func (f *fakeSupplierRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func TestCreatePR_NotApprove(t *testing.T) {
	purchaseRequestRepo := &fakePurchaseRequestRepository{
		request: &purchaserequestentity.PurchaseRequest{
			ID:          10,
			Status:      "DRAFT",
			WarehouseID: 5,
			Items: []*purchaserequestentity.PurchaseRequestItem{
				{ProductID: 7, Quantity: 3},
			},
		},
	}
	supplierRepo := &fakeSupplierRepository{supplier: &supplierentity.Supplier{ID: 4, IsActive: 1}}
	orderRepo := &fakePurchaseOrderRepository{}
	uc := NewUsecase(orderRepo, purchaseRequestRepo, supplierRepo)

	_, err := uc.CreateFromPurchaseRequest(context.Background(), 10, 4)
	if err == nil {
		t.Fatal("CreateFromPurchaseRequest() expected error when purchase request is not approved")
	}
	if err != ErrPurchaseRequestNotApproved {
		t.Fatalf("expected ErrPurchaseRequestNotApproved, got %v", err)
	}
}

func TestCreatePR_AlreadyExists(t *testing.T) {
	purchaseRequestRepo := &fakePurchaseRequestRepository{
		request: &purchaserequestentity.PurchaseRequest{
			ID:          10,
			Status:      "APPROVED",
			WarehouseID: 5,
			Items: []*purchaserequestentity.PurchaseRequestItem{
				{ProductID: 7, Quantity: 3},
			},
		},
	}
	supplierRepo := &fakeSupplierRepository{supplier: &supplierentity.Supplier{ID: 4, IsActive: 1}}
	orderRepo := &fakePurchaseOrderRepository{
		order: &purchaseorderentity.PurchaseOrder{ID: 99, PurchaseRequestID: 10},
	}
	uc := NewUsecase(orderRepo, purchaseRequestRepo, supplierRepo)

	_, err := uc.CreateFromPurchaseRequest(context.Background(), 10, 4)
	if err == nil {
		t.Fatal("CreateFromPurchaseRequest() expected error when purchase order already exists")
	}
	if err != ErrPurchaseOrderAlreadyExists {
		t.Fatalf("expected ErrPurchaseOrderAlreadyExists, got %v", err)
	}
}
