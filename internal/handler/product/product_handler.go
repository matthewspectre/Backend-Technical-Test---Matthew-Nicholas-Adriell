package product

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	entity "be_evindo/internal/entity/product"
	usecase "be_evindo/internal/usecase/product"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	uc usecase.Usecase
}

const responseTimeLayout = "2006-01-02 15:04:05"

func formatResponseTime(value string) *string {
	for _, layout := range []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999-07",
		"2006-01-02 15:04:05.999999+07",
		"2006-01-02 15:04:05-07",
	} {
		if parsed, err := time.Parse(layout, value); err == nil {
			formatted := parsed.Format(responseTimeLayout)
			return &formatted
		}
	}
	if value == "" {
		return nil
	}
	return &value
}

func productResponse(data *entity.Product) ProductResponse {
	return ProductResponse{
		ID:        data.ID,
		SKU:       data.SKU,
		Name:      data.Name,
		Unit:      data.Unit,
		IsActive:  data.IsActive,
		CreatedAt: data.CreatedAt.Format(responseTimeLayout),
		UpdatedAt: formatResponseTime(data.UpdatedAt),
	}
}

func isDuplicateSKUError(err error) bool {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return false
	}
	return postgresError.Code == "23505" &&
		(postgresError.ConstraintName == "product_unique" || postgresError.ColumnName == "sku")
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (handler *Handler) Create(context *gin.Context) {
	var request ProductCreateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := &entity.Product{
		SKU:      request.SKU,
		Name:     request.Name,
		Unit:     request.Unit,
		IsActive: request.IsActive,
	}
	if data.IsActive == 0 {
		data.IsActive = 1
	}

	if err := handler.uc.Create(context.Request.Context(), data); err != nil {
		if isDuplicateSKUError(err) {
			context.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "SKU_NOT_APPROVED",
					"message": "SKU already exists. Please provide a unique SKU.",
				},
			})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProduct, err := handler.uc.FindByID(context.Request.Context(), data.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := productResponse(createdProduct)
	response.UpdatedAt = nil
	context.JSON(http.StatusCreated, response)
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	context.JSON(http.StatusOK, productResponse(data))
}
