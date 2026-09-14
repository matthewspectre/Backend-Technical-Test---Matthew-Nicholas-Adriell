package purchase_request

import (
	"context"
	"testing"

	productentity "be_evindo/internal/entity/product"
	entity "be_evindo/internal/entity/purchase_request"
	warehouseentity "be_evindo/internal/entity/warehouse"
)

type fakePurchaseRequestRepository struct {
	created *entity.PurchaseRequest
}

func (f *fakePurchaseRequestRepository) Create(ctx context.Context, data *entity.PurchaseRequest) error {
	f.created = data
	return nil
}

func (f *fakePurchaseRequestRepository) Update(ctx context.Context, id int64, warehouseID *int, status *string, items []*entity.PurchaseRequestItem) error {
	return nil
}

func (f *fakePurchaseRequestRepository) Approve(ctx context.Context, id int64) error {
	return nil
}

func (f *fakePurchaseRequestRepository) Reject(ctx context.Context, id int64) error {
	return nil
}

func (f *fakePurchaseRequestRepository) FindAll(ctx context.Context, status string) ([]*entity.PurchaseRequest, error) {
	return nil, nil
}

func (f *fakePurchaseRequestRepository) FindByID(ctx context.Context, id int64) (*entity.PurchaseRequest, error) {
	return nil, nil
}

type fakeWarehouseRepository struct {
	warehouse *warehouseentity.Warehouse
}

func (f *fakeWarehouseRepository) FindAll(ctx context.Context) ([]*warehouseentity.Warehouse, error) {
	return nil, nil
}

func (f *fakeWarehouseRepository) FindByID(ctx context.Context, id int) (*warehouseentity.Warehouse, error) {
	if f.warehouse == nil {
		return nil, nil
	}
	return f.warehouse, nil
}

func (f *fakeWarehouseRepository) Create(ctx context.Context, entity *warehouseentity.Warehouse) error {
	return nil
}

func (f *fakeWarehouseRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	return nil
}

func (f *fakeWarehouseRepository) Delete(ctx context.Context, id int) error {
	return nil
}

type fakeProductRepository struct {
	product *productentity.Product
}

func (f *fakeProductRepository) FindAll(ctx context.Context) ([]*productentity.Product, error) {
	return nil, nil
}

func (f *fakeProductRepository) FindByID(ctx context.Context, id int) (*productentity.Product, error) {
	if f.product == nil {
		return nil, nil
	}
	return f.product, nil
}

func (f *fakeProductRepository) Create(ctx context.Context, entity *productentity.Product) error {
	return nil
}

func (f *fakeProductRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	return nil
}

func (f *fakeProductRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func TestCreate_EmptyItems_ShouldReturnError(t *testing.T) {
	warehouseRepo := &fakeWarehouseRepository{warehouse: &warehouseentity.Warehouse{ID: 1, IsActive: 1}}
	productRepo := &fakeProductRepository{product: &productentity.Product{ID: 10, IsActive: 1}}
	repo := &fakePurchaseRequestRepository{}
	uc := NewUsecase(repo, warehouseRepo, productRepo)

	data := &entity.PurchaseRequest{
		RequestNumber: "PR-2026-000004",
		WarehouseID:   1,
		RequestedBy:   2,
		Items:         []*entity.PurchaseRequestItem{},
	}

	err := uc.Create(context.Background(), data)
	if err == nil {
		t.Fatal("Create() expected error when items is empty")
	}
	if err != ErrItemsRequired {
		t.Fatalf("expected ErrItemsRequired, got %v", err)
	}
}
