package purchase_order

import (
	"context"
	"fmt"
	"time"

	entity "be_evindo/internal/entity/purchase_order"
	model "be_evindo/internal/model/purchase_order"
	repo "be_evindo/internal/repository/purchase_order"

	"gorm.io/gorm"
)

type RepositoryPostgre struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.PurchaseOrderRepository {
	return &RepositoryPostgre{db: db}
}

func toModel(data *entity.PurchaseOrder) *model.PurchaseOrderModel {
	if data == nil {
		return nil
	}
	return &model.PurchaseOrderModel{
		ID: data.ID, PONumber: data.PONumber,
		PurchaseRequestID: data.PurchaseRequestID, SupplierID: data.SupplierID,
		WarehouseID: data.WarehouseID, Status: data.Status,
		CreatedAt: data.CreatedAt, UpdatedAt: data.UpdatedAt,
	}
}

func toEntity(data *model.PurchaseOrderModel) *entity.PurchaseOrder {
	if data == nil {
		return nil
	}
	items := make([]*entity.PurchaseOrderItem, 0, len(data.Items))
	for _, item := range data.Items {
		items = append(items, &entity.PurchaseOrderItem{
			ID: item.ID, PurchaseOrderID: item.PurchaseOrderID, ProductID: item.ProductID,
			ProductName: item.ProductName, Unit: item.Unit,
			OrderedQuantity: item.OrderedQuantity, ReceivedQuantity: item.ReceivedQuantity,
		})
	}
	return &entity.PurchaseOrder{
		ID: data.ID, PONumber: data.PONumber,
		PurchaseRequestID: data.PurchaseRequestID, SupplierID: data.SupplierID,
		SupplierName: data.SupplierName, WarehouseID: data.WarehouseID,
		WarehouseName: data.WarehouseName, Status: data.Status, Items: items,
		CreatedAt: data.CreatedAt, UpdatedAt: data.UpdatedAt,
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, data *entity.PurchaseOrder, items []*entity.PurchaseOrderItem) error {
	if data.PONumber == "" {
		data.PONumber = fmt.Sprintf("PO-%s-%d", time.Now().UTC().Format("2006"), time.Now().UTC().UnixNano())
	}
	return repository.db.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		orderModel := toModel(data)
		if err := transaction.Create(orderModel).Error; err != nil {
			return err
		}
		for _, item := range items {
			itemModel := &model.PurchaseOrderItemModel{
				PurchaseOrderID: orderModel.ID, ProductID: item.ProductID,
				OrderedQuantity: item.OrderedQuantity, ReceivedQuantity: item.ReceivedQuantity,
			}
			if err := transaction.Create(itemModel).Error; err != nil {
				return err
			}
		}
		data.ID = orderModel.ID
		return nil
	})
}

func (repository *RepositoryPostgre) UpdateStatus(ctx context.Context, id int64, status string) error {
	result := repository.db.WithContext(ctx).
		Model(&model.PurchaseOrderModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repository *RepositoryPostgre) FindAll(ctx context.Context, status string) ([]*entity.PurchaseOrder, error) {
	var orders []*model.PurchaseOrderModel
	query := repository.db.WithContext(ctx).
		Table("purchase_orders po").
		Select("po.*, s.name AS supplier_name, w.name AS warehouse_name").
		Joins("LEFT JOIN supplier s ON s.id = po.supplier_id").
		Joins("LEFT JOIN warehouse w ON w.id = po.warehouse_id")
	if status != "" {
		query = query.Where("po.status = ?", status)
	}
	if err := query.Order("po.id DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	if err := repository.loadItems(ctx, orders); err != nil {
		return nil, err
	}
	result := make([]*entity.PurchaseOrder, 0, len(orders))
	for _, order := range orders {
		result = append(result, toEntity(order))
	}
	return result, nil
}

func (repository *RepositoryPostgre) FindByID(ctx context.Context, id int64) (*entity.PurchaseOrder, error) {
	orderModel := &model.PurchaseOrderModel{}
	if err := repository.db.WithContext(ctx).
		Table("purchase_orders po").
		Select("po.*, s.name AS supplier_name, w.name AS warehouse_name").
		Joins("LEFT JOIN supplier s ON s.id = po.supplier_id").
		Joins("LEFT JOIN warehouse w ON w.id = po.warehouse_id").
		Where("po.id = ?", id).First(orderModel).Error; err != nil {
		return nil, err
	}
	if err := repository.loadItems(ctx, []*model.PurchaseOrderModel{orderModel}); err != nil {
		return nil, err
	}
	return toEntity(orderModel), nil
}

func (repository *RepositoryPostgre) loadItems(ctx context.Context, orders []*model.PurchaseOrderModel) error {
	for _, order := range orders {
		var items []*model.PurchaseOrderItemModel
		if err := repository.db.WithContext(ctx).
			Table("purchase_order_items poi").
			Select("poi.*, p.name AS product_name, p.unit AS unit").
			Joins("LEFT JOIN product p ON p.id = poi.product_id").
			Where("poi.purchase_order_id = ?", order.ID).
			Order("poi.id ASC").Find(&items).Error; err != nil {
			return err
		}
		order.Items = items
	}
	return nil
}

func (repository *RepositoryPostgre) FindByPurchaseRequestID(ctx context.Context, purchaseRequestID int64) (*entity.PurchaseOrder, error) {
	orderModel := &model.PurchaseOrderModel{}
	if err := repository.db.WithContext(ctx).
		Table("purchase_orders po").
		Select("po.*, s.name AS supplier_name, w.name AS warehouse_name").
		Joins("LEFT JOIN supplier s ON s.id = po.supplier_id").
		Joins("LEFT JOIN warehouse w ON w.id = po.warehouse_id").
		Where("po.purchase_request_id = ?", purchaseRequestID).First(orderModel).Error; err != nil {
		return nil, err
	}
	if err := repository.loadItems(ctx, []*model.PurchaseOrderModel{orderModel}); err != nil {
		return nil, err
	}
	return toEntity(orderModel), nil
}
