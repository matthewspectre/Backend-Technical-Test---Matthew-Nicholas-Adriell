package inventory_movement

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	entity "be_evindo/internal/entity/inventory_movement"
	usecase "be_evindo/internal/usecase/inventory_movement"

	"github.com/gin-gonic/gin"
)

type Handler struct{ uc usecase.Usecase }

func NewHandler(uc usecase.Usecase) *Handler { return &Handler{uc: uc} }

func (handler *Handler) GetAll(context *gin.Context) {
	var productID *int
	if context.Query("product_id") != "" {
		parsedProductID, err := strconv.Atoi(context.Query("product_id"))
		if err != nil || parsedProductID <= 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
			return
		}
		productID = &parsedProductID
	}
	var warehouseID *int
	if context.Query("warehouse_id") != "" {
		parsedWarehouseID, err := strconv.Atoi(context.Query("warehouse_id"))
		if err != nil || parsedWarehouseID <= 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse_id"})
			return
		}
		warehouseID = &parsedWarehouseID
	}
	reference := strings.TrimSpace(context.Query("reference"))

	data, err := handler.uc.FindAll(context.Request.Context(), productID, warehouseID, reference)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := make([]InventoryMovementResponse, 0, len(data))
	for _, movement := range data {
		response = append(response, inventoryMovementResponse(movement))
	}
	context.JSON(http.StatusOK, InventoryMovementListResponse{Data: response, TotalData: len(response)})
}

func inventoryMovementResponse(data *entity.InventoryMovement) InventoryMovementResponse {
	return InventoryMovementResponse{
		ID:           data.ID,
		WarehouseID:  data.WarehouseID,
		ProductID:    data.ProductID,
		MovementType: data.MovementType,
		Quantity:     data.Quantity,
		Reference:    data.Reference,
		CreatedAt:    data.CreatedAt.Format(time.RFC3339),
	}
}
