package goods_receipt

import (
	"context"
	"errors"

	entity "be_evindo/internal/entity/goods_receipt"
	purchaseorder "be_evindo/internal/entity/purchase_order"
	repo "be_evindo/internal/repository/goods_receipt"
	purchaseorderrepo "be_evindo/internal/repository/purchase_order"
)

type Usecase interface {
	Create(ctx context.Context, data *entity.GoodsReceipt) error
	FindAll(ctx context.Context) ([]*entity.GoodsReceipt, error)
	FindByID(ctx context.Context, id int64) (*entity.GoodsReceipt, error)
}

var ErrPurchaseOrderNotReceivable = errors.New("purchase order must be ORDERED or PARTIALLY_RECEIVED")
var ErrPurchaseOrderAlreadyReceived = errors.New("purchase order has already been fully received")
var ErrPurchaseOrderCancelled = errors.New("purchase order has been cancelled and cannot receive goods")
var ErrReceiptItemsRequired = errors.New("at least one goods receipt item is required")
var ErrReceiptProductNotInOrder = errors.New("product received must be a product that exists in the purchase order")
var ErrReceiptQuantityInvalid = errors.New("received quantity must be greater than zero")
var ErrReceiptQuantityExceeded = errors.New("received quantity exceeds ordered quantity")

type purchaseOrderReader interface {
	FindByID(ctx context.Context, id int64) (*purchaseorder.PurchaseOrder, error)
}

type usecase struct {
	repository     repo.GoodsReceiptRepository
	purchaseOrders purchaseOrderReader
}

func NewUsecase(repository repo.GoodsReceiptRepository, purchaseOrders purchaseorderrepo.PurchaseOrderRepository) Usecase {
	return &usecase{repository: repository, purchaseOrders: purchaseOrders}
}

func (usecase *usecase) Create(ctx context.Context, data *entity.GoodsReceipt) error {
	if data == nil {
		return errors.New("goods receipt is required")
	}
	if data.ReceivedBy <= 0 {
		return errors.New("received_by must be greater than zero")
	}
	if data.PurchaseOrderID <= 0 {
		return errors.New("purchase_order_id must be greater than zero")
	}
	if len(data.Items) == 0 {
		return ErrReceiptItemsRequired
	}
	purchaseOrder, err := usecase.purchaseOrders.FindByID(ctx, data.PurchaseOrderID)
	if err != nil {
		return err
	}
	switch purchaseOrder.Status {
	case "ORDERED", "PARTIALLY_RECEIVED":
		// allowed
	case "RECEIVED":
		return ErrPurchaseOrderAlreadyReceived
	case "CANCELLED":
		return ErrPurchaseOrderCancelled
	default:
		return ErrPurchaseOrderNotReceivable
	}
	itemsByProduct := make(map[int]*purchaseorder.PurchaseOrderItem, len(purchaseOrder.Items))
	for _, item := range purchaseOrder.Items {
		itemsByProduct[item.ProductID] = item
	}
	seen := make(map[int]struct{}, len(data.Items))
	for _, item := range data.Items {
		orderItem, ok := itemsByProduct[item.ProductID]
		if !ok {
			return ErrReceiptProductNotInOrder
		}
		item.PurchaseOrderItemID = orderItem.ID
		if item.ReceivedQuantity <= 0 {
			return ErrReceiptQuantityInvalid
		}
		if orderItem.ReceivedQuantity+item.ReceivedQuantity > orderItem.OrderedQuantity {
			return ErrReceiptQuantityExceeded
		}
		if _, exists := seen[item.ProductID]; exists {
			return errors.New("a product may only appear once per goods receipt")
		}
		seen[item.ProductID] = struct{}{}
	}
	data.Status = "POSTED"
	data.WarehouseID = purchaseOrder.WarehouseID
	err = usecase.repository.Create(ctx, data)
	if errors.Is(err, repo.ErrProductNotInPurchaseOrder) {
		return ErrReceiptProductNotInOrder
	}
	if errors.Is(err, repo.ErrReceivedQuantityExceeded) {
		return ErrReceiptQuantityExceeded
	}
	return err
}

func (usecase *usecase) FindAll(ctx context.Context) ([]*entity.GoodsReceipt, error) {
	return usecase.repository.FindAll(ctx)
}

func (usecase *usecase) FindByID(ctx context.Context, id int64) (*entity.GoodsReceipt, error) {
	if id <= 0 {
		return nil, errors.New("goods receipt id must be greater than zero")
	}
	return usecase.repository.FindByID(ctx, id)
}
