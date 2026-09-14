package purchase_order

import (
	"errors"
	"net/http"
	"strconv"

	entity "be_evindo/internal/entity/purchase_order"
	usecase "be_evindo/internal/usecase/purchase_order"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	uc usecase.Usecase
}

const responseTimeLayout = "2006-01-02 15:04:05"

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func purchaseOrderResponse(data *entity.PurchaseOrder) PurchaseOrderResponse {
	items := make([]PurchaseOrderItemResponse, 0, len(data.Items))
	for _, item := range data.Items {
		items = append(items, PurchaseOrderItemResponse{
			ID: item.ID, ProductID: item.ProductID, ProductName: item.ProductName,
			Unit: item.Unit, OrderedQuantity: item.OrderedQuantity,
			ReceivedQuantity: item.ReceivedQuantity,
		})
	}
	return PurchaseOrderResponse{
		ID: data.ID, PONumber: data.PONumber, PurchaseRequestID: data.PurchaseRequestID,
		SupplierID: data.SupplierID, SupplierName: data.SupplierName,
		WarehouseID: data.WarehouseID, WarehouseName: data.WarehouseName,
		Status: data.Status, Items: items, CreatedAt: data.CreatedAt.Format(responseTimeLayout),
		UpdatedAt: data.UpdatedAt.Format(responseTimeLayout),
	}
}

func (handler *Handler) CreateFromPurchaseRequest(context *gin.Context) {
	purchaseRequestID, err := strconv.ParseInt(context.Param("purchase_request_id"), 10, 64)
	if err != nil || purchaseRequestID <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase request id"})
		return
	}
	var request PurchaseOrderCreateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := handler.uc.CreateFromPurchaseRequest(context.Request.Context(), purchaseRequestID, request.SupplierID)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusCreated, purchaseOrderResponse(data))
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase order id"})
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, purchaseOrderResponse(data))
}

func (handler *Handler) UpdateStatus(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase order id"})
		return
	}
	var request PurchaseOrderStatusUpdateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := handler.uc.UpdateStatus(context.Request.Context(), id, request.Status); err != nil {
		handler.writeError(context, err)
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, purchaseOrderResponse(data))
}

func (handler *Handler) writeError(context *gin.Context, err error) {
	if errors.Is(err, usecase.ErrInvalidPurchaseOrderStatus) {
		context.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "PURCHASE_ORDER_STATUS_INVALID", "message": "Status must be DRAFT, ORDERED, PARTIALLY_RECEIVED, RECEIVED, or CANCELLED.",
		}})
		return
	}
	if errors.Is(err, usecase.ErrPurchaseRequestNotApproved) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code": "PURCHASE_REQUEST_NOT_APPROVED", "message": "Only APPROVED purchase requests can create purchase orders.",
		}})
		return
	}
	if errors.Is(err, usecase.ErrPurchaseOrderAlreadyExists) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code": "PURCHASE_ORDER_ALREADY_EXISTS", "message": "A purchase order already exists for this purchase request.",
		}})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "purchase order or related data not found"})
		return
	}
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
