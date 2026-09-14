package goods_receipt

import (
	"context"
	"fmt"
	"time"

	entity "be_evindo/internal/entity/goods_receipt"
	model "be_evindo/internal/model/goods_receipt"
	pomodel "be_evindo/internal/model/purchase_order"
	repo "be_evindo/internal/repository/goods_receipt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryPostgre struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) repo.GoodsReceiptRepository { return &RepositoryPostgre{db: db} }

func toModel(data *entity.GoodsReceipt) *model.GoodsReceiptModel {
	items := make([]*model.GoodsReceiptItemModel, 0, len(data.Items))
	for _, item := range data.Items {
		items = append(items, &model.GoodsReceiptItemModel{
			GoodsReceiptID: item.GoodsReceiptID, PurchaseOrderItemID: item.PurchaseOrderItemID,
			ProductID: item.ProductID, ReceivedQuantity: item.ReceivedQuantity,
		})
	}
	return &model.GoodsReceiptModel{
		ID: data.ID, ReceiptNumber: data.ReceiptNumber, PurchaseOrderID: data.PurchaseOrderID,
		WarehouseID: data.WarehouseID, ReceivedBy: data.ReceivedBy, Status: data.Status,
		ReceivedAt: data.ReceivedAt, CreatedAt: data.CreatedAt, UpdatedAt: data.UpdatedAt,
		Items: items,
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, data *entity.GoodsReceipt) error {
	if data.ReceiptNumber == "" {
		data.ReceiptNumber = fmt.Sprintf("GR-%s-%d", time.Now().UTC().Format("2006"), time.Now().UTC().UnixNano())
	}
	return repository.db.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		var purchaseOrder pomodel.PurchaseOrderModel
		if err := transaction.Clauses(clauseLocking()).First(&purchaseOrder, data.PurchaseOrderID).Error; err != nil {
			return err
		}
		var orderItems []*pomodel.PurchaseOrderItemModel
		if err := transaction.Clauses(clauseLocking()).Where("purchase_order_id = ?", data.PurchaseOrderID).Find(&orderItems).Error; err != nil {
			return err
		}
		itemsByID := make(map[int64]*pomodel.PurchaseOrderItemModel, len(orderItems))
		for _, item := range orderItems {
			itemsByID[item.ID] = item
		}
		for _, receiptItem := range data.Items {
			orderItem, ok := itemsByID[receiptItem.PurchaseOrderItemID]
			if !ok || orderItem.ProductID != receiptItem.ProductID {
				return repo.ErrProductNotInPurchaseOrder
			}
			if receiptItem.ReceivedQuantity <= 0 || orderItem.ReceivedQuantity+receiptItem.ReceivedQuantity > orderItem.OrderedQuantity {
				return repo.ErrReceivedQuantityExceeded
			}
		}
		receiptModel := toModel(data)
		if err := transaction.Create(receiptModel).Error; err != nil {
			return err
		}
		data.ID = receiptModel.ID
		for _, receiptItem := range receiptModel.Items {
			receiptItem.GoodsReceiptID = receiptModel.ID
			if err := transaction.Create(receiptItem).Error; err != nil {
				return err
			}
			if err := transaction.Model(&pomodel.PurchaseOrderItemModel{}).Where("id = ?", receiptItem.PurchaseOrderItemID).UpdateColumn("received_quantity", gorm.Expr("received_quantity + ?", receiptItem.ReceivedQuantity)).Error; err != nil {
				return err
			}
		}
		var remaining int64
		if err := transaction.Model(&pomodel.PurchaseOrderItemModel{}).Where("purchase_order_id = ? AND received_quantity < ordered_quantity", data.PurchaseOrderID).Count(&remaining).Error; err != nil {
			return err
		}
		status := "RECEIVED"
		if remaining > 0 {
			status = "PARTIALLY_RECEIVED"
		}
		return transaction.Model(&pomodel.PurchaseOrderModel{}).Where("id = ?", data.PurchaseOrderID).Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
	})
}

func toEntity(data *model.GoodsReceiptModel) *entity.GoodsReceipt {
	items := make([]*entity.GoodsReceiptItem, 0, len(data.Items))
	for _, item := range data.Items {
		items = append(items, &entity.GoodsReceiptItem{
			ID: item.ID, GoodsReceiptID: item.GoodsReceiptID,
			PurchaseOrderItemID: item.PurchaseOrderItemID,
			ProductID:           item.ProductID, ReceivedQuantity: item.ReceivedQuantity,
		})
	}
	return &entity.GoodsReceipt{
		ID: data.ID, ReceiptNumber: data.ReceiptNumber, PurchaseOrderID: data.PurchaseOrderID,
		WarehouseID: data.WarehouseID, ReceivedBy: data.ReceivedBy, Status: data.Status,
		ReceivedAt: data.ReceivedAt, CreatedAt: data.CreatedAt, UpdatedAt: data.UpdatedAt,
		Items: items,
	}
}

func (repository *RepositoryPostgre) FindAll(ctx context.Context) ([]*entity.GoodsReceipt, error) {
	var receipts []*model.GoodsReceiptModel
	if err := repository.db.WithContext(ctx).Order("id DESC").Find(&receipts).Error; err != nil {
		return nil, err
	}
	if err := repository.loadItems(ctx, receipts); err != nil {
		return nil, err
	}
	result := make([]*entity.GoodsReceipt, 0, len(receipts))
	for _, receipt := range receipts {
		result = append(result, toEntity(receipt))
	}
	return result, nil
}

func (repository *RepositoryPostgre) FindByID(ctx context.Context, id int64) (*entity.GoodsReceipt, error) {
	receipt := &model.GoodsReceiptModel{}
	if err := repository.db.WithContext(ctx).First(receipt, id).Error; err != nil {
		return nil, err
	}
	if err := repository.loadItems(ctx, []*model.GoodsReceiptModel{receipt}); err != nil {
		return nil, err
	}
	return toEntity(receipt), nil
}

func (repository *RepositoryPostgre) loadItems(ctx context.Context, receipts []*model.GoodsReceiptModel) error {
	for _, receipt := range receipts {
		var items []*model.GoodsReceiptItemModel
		if err := repository.db.WithContext(ctx).
			Table("goods_receipt_items gri").
			Select("gri.*").
			Where("gri.goods_receipt_id = ?", receipt.ID).
			Order("gri.id ASC").Find(&items).Error; err != nil {
			return err
		}
		receipt.Items = items
	}
	return nil
}

func clauseLocking() clause.Locking { return clause.Locking{Strength: "UPDATE"} }
