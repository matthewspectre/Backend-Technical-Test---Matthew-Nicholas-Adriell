package goods_receipt

import (
	"errors"
	"net/http"
	"strconv"

	entity "be_evindo/internal/entity/goods_receipt"
	usecase "be_evindo/internal/usecase/goods_receipt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ uc usecase.Usecase }

func NewHandler(uc usecase.Usecase) *Handler { return &Handler{uc: uc} }

func (handler *Handler) Create(context *gin.Context) {
	userIDValue, exists := context.Get("user_id")
	userID, ok := userIDValue.(int64)
	if !exists || !ok || userID <= 0 {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user identity"})
		return
	}
	purchaseOrderID, err := strconv.ParseInt(context.Param("purchase_order_id"), 10, 64)
	if err != nil || purchaseOrderID <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase order id"})
		return
	}
	var request GoodsReceiptCreateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := &entity.GoodsReceipt{
		PurchaseOrderID: purchaseOrderID,
		ReceivedBy:      userID,
		Items:           make([]*entity.GoodsReceiptItem, 0, len(request.Items)),
	}
	for _, item := range request.Items {
		data.Items = append(data.Items, &entity.GoodsReceiptItem{
			ProductID: item.ProductID, ReceivedQuantity: item.ReceivedQuantity,
		})
	}
	if err := handler.uc.Create(context.Request.Context(), data); err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusCreated, GoodsReceiptResponse{
		ID: data.ID, ReceiptNumber: data.ReceiptNumber, PurchaseOrderID: data.PurchaseOrderID,
		WarehouseID: data.WarehouseID, ReceivedBy: data.ReceivedBy, Status: data.Status,
	})
}

func (handler *Handler) GetAll(context *gin.Context) {
	data, err := handler.uc.FindAll(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	response := make([]GoodsReceiptResponse, 0, len(data))
	for _, receipt := range data {
		response = append(response, goodsReceiptResponse(receipt))
	}
	context.JSON(http.StatusOK, GoodsReceiptListResponse{Data: response, TotalData: len(response)})
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid goods receipt id"})
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, goodsReceiptResponse(data))
}

func goodsReceiptResponse(data *entity.GoodsReceipt) GoodsReceiptResponse {
	items := make([]GoodsReceiptItemResponse, 0, len(data.Items))
	for _, item := range data.Items {
		items = append(items, GoodsReceiptItemResponse{
			ID: item.ID, PurchaseOrderItemID: item.PurchaseOrderItemID,
			ProductID: item.ProductID, ReceivedQuantity: item.ReceivedQuantity,
		})
	}
	return GoodsReceiptResponse{
		ID: data.ID, ReceiptNumber: data.ReceiptNumber, PurchaseOrderID: data.PurchaseOrderID,
		WarehouseID: data.WarehouseID, ReceivedBy: data.ReceivedBy, Status: data.Status, Items: items,
	}
}

func (handler *Handler) writeError(context *gin.Context, err error) {
	status := http.StatusBadRequest
	code := "GOODS_RECEIPT_INVALID"
	message := err.Error()
	switch {
	case errors.Is(err, usecase.ErrPurchaseOrderNotReceivable):
		status, code, message = http.StatusConflict, "PURCHASE_ORDER_NOT_RECEIVABLE", "Purchase Order must be ORDERED or PARTIALLY_RECEIVED."
	case errors.Is(err, usecase.ErrReceiptItemsRequired):
		code, message = "GOODS_RECEIPT_ITEMS_REQUIRED", "At least one goods receipt item is required."
	case errors.Is(err, usecase.ErrReceiptProductNotInOrder):
		code, message = "GOODS_RECEIPT_PRODUCT_NOT_IN_ORDER", "Product is not part of the Purchase Order."
	case errors.Is(err, usecase.ErrReceiptQuantityExceeded):
		status, code, message = http.StatusConflict, "GOODS_RECEIPT_QUANTITY_EXCEEDED", "Received quantity exceeds ordered quantity."
	case errors.Is(err, gorm.ErrRecordNotFound):
		status, code, message = http.StatusNotFound, "PURCHASE_ORDER_NOT_FOUND", "Purchase Order not found."
	}
	context.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}
