package goods_receipt

import (
	"context"
	"errors"

	entity "be_evindo/internal/entity/goods_receipt"
)

var ErrProductNotInPurchaseOrder = errors.New("product is not part of the purchase order")
var ErrReceivedQuantityExceeded = errors.New("received quantity exceeds ordered quantity")

type GoodsReceiptRepository interface {
	Create(ctx context.Context, data *entity.GoodsReceipt) error
	FindAll(ctx context.Context) ([]*entity.GoodsReceipt, error)
	FindByID(ctx context.Context, id int64) (*entity.GoodsReceipt, error)
}
