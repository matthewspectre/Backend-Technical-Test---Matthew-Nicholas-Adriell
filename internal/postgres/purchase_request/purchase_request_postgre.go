package purchase_request

import (
	"context"
	"time"

	entity "be_evindo/internal/entity/purchase_request"
	model "be_evindo/internal/model/purchase_request"
	repo "be_evindo/internal/repository/purchase_request"

	"gorm.io/gorm"
)

type RepositoryPostgre struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.PurchaseRequestRepository {
	return &RepositoryPostgre{db: db}
}

func toModel(data *entity.PurchaseRequest) *model.PurchaseRequestModel {
	if data == nil {
		return nil
	}
	items := make([]*model.PurchaseRequestItemModel, 0, len(data.Items))
	for _, item := range data.Items {
		items = append(items, &model.PurchaseRequestItemModel{
			ID: item.ID, PurchaseRequestID: item.PurchaseRequestID,
			ProductID: item.ProductID, Quantity: item.Quantity,
		})
	}
	return &model.PurchaseRequestModel{
		ID: data.ID, RequestNumber: data.RequestNumber, WarehouseID: data.WarehouseID,
		RequestedBy: data.RequestedBy, Status: data.Status,
		Items:     items,
		CreatedAt: data.CreatedAt, UpdatedAt: data.UpdatedAt,
	}
}

func toEntity(data *model.PurchaseRequestModel) *entity.PurchaseRequest {
	if data == nil {
		return nil
	}
	items := make([]*entity.PurchaseRequestItem, 0, len(data.Items))
	for _, item := range data.Items {
		items = append(items, &entity.PurchaseRequestItem{
			ID: item.ID, PurchaseRequestID: item.PurchaseRequestID,
			ProductID: item.ProductID, ProductName: item.ProductName,
			Unit: item.Unit, Quantity: item.Quantity,
		})
	}
	return &entity.PurchaseRequest{
		ID: data.ID, RequestNumber: data.RequestNumber, WarehouseID: data.WarehouseID,
		WarehouseName: data.WarehouseName, RequestedBy: data.RequestedBy,
		RequesterName: data.RequesterName, Status: data.Status,
		Items:     items,
		CreatedAt: data.CreatedAt, UpdatedAt: data.UpdatedAt,
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, data *entity.PurchaseRequest) error {
	return repository.db.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		requestModel := toModel(data)
		items := requestModel.Items
		requestModel.Items = nil
		if err := transaction.Create(requestModel).Error; err != nil {
			return err
		}
		for _, item := range items {
			item.PurchaseRequestID = requestModel.ID
			if err := transaction.Create(item).Error; err != nil {
				return err
			}
		}
		if data != nil {
			data.ID = requestModel.ID
		}
		return nil
	})
}

func (repository *RepositoryPostgre) Update(ctx context.Context, id int64, warehouseID *int, status *string, items []*entity.PurchaseRequestItem) error {
	return repository.db.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		updates := map[string]interface{}{"updated_at": time.Now()}
		if warehouseID != nil {
			updates["warehouse_id"] = *warehouseID
		}
		if status != nil {
			updates["status"] = *status
		}
		result := transaction.Model(&model.PurchaseRequestModel{}).
			Where("id = ?", id).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if items != nil {
			if err := transaction.Where("purchase_request_id = ?", id).
				Delete(&model.PurchaseRequestItemModel{}).Error; err != nil {
				return err
			}
			for _, item := range items {
				itemModel := &model.PurchaseRequestItemModel{
					PurchaseRequestID: id,
					ProductID:         item.ProductID,
					Quantity:          item.Quantity,
				}
				if err := transaction.Create(itemModel).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (repository *RepositoryPostgre) Approve(ctx context.Context, id int64) error {
	result := repository.db.WithContext(ctx).
		Model(&model.PurchaseRequestModel{}).
		Where("id = ? AND status = ?", id, "SUBMITTED").
		Updates(map[string]interface{}{"status": "APPROVED", "updated_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repository *RepositoryPostgre) Reject(ctx context.Context, id int64) error {
	result := repository.db.WithContext(ctx).
		Model(&model.PurchaseRequestModel{}).
		Where("id = ? AND status = ?", id, "SUBMITTED").
		Updates(map[string]interface{}{"status": "REJECTED", "updated_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repository *RepositoryPostgre) FindAll(ctx context.Context, status string) ([]*entity.PurchaseRequest, error) {
	var models []*model.PurchaseRequestModel
	query := repository.db.WithContext(ctx).
		Table("purchase_requests pr").
		Select("pr.*, w.name AS warehouse_name, u.name AS requester_name").
		Joins("LEFT JOIN warehouse w ON w.id = pr.warehouse_id").
		Joins("LEFT JOIN users u ON u.id = pr.requested_by")
	if status != "" {
		query = query.Where("pr.status = ?", status)
	}
	if err := query.Order("pr.id DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	if err := repository.loadItems(ctx, models); err != nil {
		return nil, err
	}
	result := make([]*entity.PurchaseRequest, 0, len(models))
	for _, request := range models {
		result = append(result, toEntity(request))
	}
	return result, nil
}

func (repository *RepositoryPostgre) FindByID(ctx context.Context, id int64) (*entity.PurchaseRequest, error) {
	requestModel := &model.PurchaseRequestModel{}
	if err := repository.db.WithContext(ctx).
		Table("purchase_requests pr").
		Select("pr.*, w.name AS warehouse_name, u.name AS requester_name").
		Joins("LEFT JOIN warehouse w ON w.id = pr.warehouse_id").
		Joins("LEFT JOIN users u ON u.id = pr.requested_by").
		Where("pr.id = ?", id).First(requestModel).Error; err != nil {
		return nil, err
	}
	if err := repository.loadItems(ctx, []*model.PurchaseRequestModel{requestModel}); err != nil {
		return nil, err
	}
	return toEntity(requestModel), nil
}

func (repository *RepositoryPostgre) loadItems(ctx context.Context, requests []*model.PurchaseRequestModel) error {
	for _, request := range requests {
		var items []*model.PurchaseRequestItemModel
		if err := repository.db.WithContext(ctx).
			Table("purchase_request_items pri").
			Select("pri.*, p.name AS product_name, p.unit AS unit").
			Joins("LEFT JOIN product p ON p.id = pri.product_id").
			Where("pri.purchase_request_id = ?", request.ID).
			Order("pri.id ASC").Find(&items).Error; err != nil {
			return err
		}
		request.Items = items
	}
	return nil
}
