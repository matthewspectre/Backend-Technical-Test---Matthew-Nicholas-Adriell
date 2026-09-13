package purchase_request

import (
	"context"

	entity "be_evindo/internal/entity/purchase_request"
)

type PurchaseRequestRepository interface {
	Create(ctx context.Context, data *entity.PurchaseRequest) error
	Update(ctx context.Context, id int64, warehouseID *int, status *string, items []*entity.PurchaseRequestItem) error
	Approve(ctx context.Context, id int64) error
	Reject(ctx context.Context, id int64) error
	FindAll(ctx context.Context, status string) ([]*entity.PurchaseRequest, error)
	FindByID(ctx context.Context, id int64) (*entity.PurchaseRequest, error)
}
