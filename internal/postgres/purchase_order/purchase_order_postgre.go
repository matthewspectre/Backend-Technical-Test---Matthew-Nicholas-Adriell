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
	return &entity.PurchaseOrder{
		ID: data.ID, PONumber: data.PONumber,
		PurchaseRequestID: data.PurchaseRequestID, SupplierID: data.SupplierID,
		SupplierName: data.SupplierName, WarehouseID: data.WarehouseID,
		WarehouseName: data.WarehouseName, Status: data.Status,
		CreatedAt: data.CreatedAt, UpdatedAt: data.UpdatedAt,
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, data *entity.PurchaseOrder) error {
	if data.PONumber == "" {
		data.PONumber = fmt.Sprintf("PO-%s-%d", time.Now().UTC().Format("2006"), time.Now().UTC().UnixNano())
	}
	orderModel := toModel(data)
	if err := repository.db.WithContext(ctx).Create(orderModel).Error; err != nil {
		return err
	}
	data.ID = orderModel.ID
	return nil
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
	return toEntity(orderModel), nil
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
	return toEntity(orderModel), nil
}
