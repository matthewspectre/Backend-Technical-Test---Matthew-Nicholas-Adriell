package goods_receipt

import (
	"context"
	"fmt"
	"testing"

	goodsreceipentity "be_evindo/internal/entity/goods_receipt"
	purchaseorderentity "be_evindo/internal/entity/purchase_order"
)

type fakeGoodsReceiptRepository struct {
	inventory map[string]int
}

func (f *fakeGoodsReceiptRepository) Create(ctx context.Context, data *goodsreceipentity.GoodsReceipt) error {
	if f.inventory == nil {
		f.inventory = make(map[string]int)
	}
	for _, item := range data.Items {
		key := fmt.Sprintf("%d:%d", data.WarehouseID, item.ProductID)
		f.inventory[key] += item.ReceivedQuantity
	}
	return nil
}

func (f *fakeGoodsReceiptRepository) FindAll(ctx context.Context) ([]*goodsreceipentity.GoodsReceipt, error) {
	return nil, nil
}

func (f *fakeGoodsReceiptRepository) FindByID(ctx context.Context, id int64) (*goodsreceipentity.GoodsReceipt, error) {
	return nil, nil
}

type fakePurchaseOrderRepository struct {
	order *purchaseorderentity.PurchaseOrder
}

func (f *fakePurchaseOrderRepository) Create(ctx context.Context, data *purchaseorderentity.PurchaseOrder, items []*purchaseorderentity.PurchaseOrderItem) error {
	return nil
}

func (f *fakePurchaseOrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return nil
}

func (f *fakePurchaseOrderRepository) FindAll(ctx context.Context, status string) ([]*purchaseorderentity.PurchaseOrder, error) {
	return nil, nil
}

func (f *fakePurchaseOrderRepository) FindByID(ctx context.Context, id int64) (*purchaseorderentity.PurchaseOrder, error) {
	if f.order == nil {
		return nil, nil
	}
	return f.order, nil
}

func (f *fakePurchaseOrderRepository) FindByPurchaseRequestID(ctx context.Context, purchaseRequestID int64) (*purchaseorderentity.PurchaseOrder, error) {
	return nil, nil
}

func TestCreate_ReceivedQuantityZeroOrLess_ShouldReturnError(t *testing.T) {
	repo := &fakeGoodsReceiptRepository{}
	poRepo := &fakePurchaseOrderRepository{
		order: &purchaseorderentity.PurchaseOrder{
			ID:          11,
			WarehouseID: 2,
			Status:      "ORDERED",
			Items: []*purchaseorderentity.PurchaseOrderItem{
				{ID: 21, ProductID: 7, OrderedQuantity: 10, ReceivedQuantity: 0},
			},
		},
	}
	uc := NewUsecase(repo, poRepo)

	data := &goodsreceipentity.GoodsReceipt{
		ReceivedBy:      1,
		PurchaseOrderID: 11,
		Items: []*goodsreceipentity.GoodsReceiptItem{
			{ProductID: 7, ReceivedQuantity: 0},
		},
	}

	err := uc.Create(context.Background(), data)
	if err == nil {
		t.Fatal("Create() expected error when received quantity is zero or less")
	}
	if err != ErrReceiptQuantityInvalid {
		t.Fatalf("expected ErrReceiptQuantityInvalid, got %v", err)
	}
}

func TestCreate_IncreasesInventoryStockByReceivedQuantity(t *testing.T) {
	repo := &fakeGoodsReceiptRepository{inventory: map[string]int{}}
	poRepo := &fakePurchaseOrderRepository{
		order: &purchaseorderentity.PurchaseOrder{
			ID:          11,
			WarehouseID: 2,
			Status:      "ORDERED",
			Items: []*purchaseorderentity.PurchaseOrderItem{
				{ID: 21, ProductID: 7, OrderedQuantity: 10, ReceivedQuantity: 0},
			},
		},
	}
	uc := NewUsecase(repo, poRepo)

	data := &goodsreceipentity.GoodsReceipt{
		ReceivedBy:      1,
		PurchaseOrderID: 11,
		WarehouseID:     2,
		Items: []*goodsreceipentity.GoodsReceiptItem{
			{ProductID: 7, ReceivedQuantity: 4},
		},
	}

	if err := uc.Create(context.Background(), data); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	got := repo.inventory["2:7"]
	if got != 4 {
		t.Fatalf("expected inventory stock to increase by 4, got %d", got)
	}
}

func TestCreate_WhenPurchaseOrderAlreadyReceived_ShouldReturnError(t *testing.T) {
	repo := &fakeGoodsReceiptRepository{inventory: map[string]int{}}
	poRepo := &fakePurchaseOrderRepository{
		order: &purchaseorderentity.PurchaseOrder{
			ID:          44,
			WarehouseID: 2,
			Status:      "RECEIVED",
			Items: []*purchaseorderentity.PurchaseOrderItem{
				{ID: 99, ProductID: 7, OrderedQuantity: 10, ReceivedQuantity: 10},
			},
		},
	}
	uc := NewUsecase(repo, poRepo)

	data := &goodsreceipentity.GoodsReceipt{
		ReceivedBy:      1,
		PurchaseOrderID: 44,
		Items: []*goodsreceipentity.GoodsReceiptItem{
			{ProductID: 7, ReceivedQuantity: 1},
		},
	}

	err := uc.Create(context.Background(), data)
	if err == nil {
		t.Fatal("Create() expected error when purchase order status is already RECEIVED")
	}
	if err != ErrPurchaseOrderAlreadyReceived {
		t.Fatalf("expected ErrPurchaseOrderAlreadyReceived, got %v", err)
	}
}
