package purchase_order

import (
	"context"

	entity "be_evindo/internal/entity/purchase_order"
)

type PurchaseOrderRepository interface {
	Create(ctx context.Context, data *entity.PurchaseOrder, items []*entity.PurchaseOrderItem) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	FindAll(ctx context.Context, status string) ([]*entity.PurchaseOrder, error)
	FindByID(ctx context.Context, id int64) (*entity.PurchaseOrder, error)
	FindByPurchaseRequestID(ctx context.Context, purchaseRequestID int64) (*entity.PurchaseOrder, error)
}
