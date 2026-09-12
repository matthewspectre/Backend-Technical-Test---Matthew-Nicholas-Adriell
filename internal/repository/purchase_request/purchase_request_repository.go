package purchase_request

import (
	"context"

	entity "be_evindo/internal/entity/purchase_request"
)

type PurchaseRequestRepository interface {
	Create(ctx context.Context, data *entity.PurchaseRequest) error
	FindAll(ctx context.Context) ([]*entity.PurchaseRequest, error)
	FindByID(ctx context.Context, id int64) (*entity.PurchaseRequest, error)
}
