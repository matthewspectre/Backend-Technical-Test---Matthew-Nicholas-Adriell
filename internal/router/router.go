package router

import (
	hgoodsreceipt "be_evindo/internal/handler/goods_receipt"
	hinventory "be_evindo/internal/handler/inventory"
	inventorymovement "be_evindo/internal/handler/inventory_movement"
	hproduct "be_evindo/internal/handler/product"
	hpurchaseorder "be_evindo/internal/handler/purchase_order"
	hpurchaserequest "be_evindo/internal/handler/purchase_request"
	hsupplier "be_evindo/internal/handler/supplier"
	userhandler "be_evindo/internal/handler/user"
	hwarehouse "be_evindo/internal/handler/warehouse"
	"be_evindo/internal/middleware"
	userusecase "be_evindo/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, userHandler *userhandler.UserHandler) {
	group := r.Group("/auth")
	group.POST("/register", userHandler.Register)
}

func RegisterLoginRoute(r *gin.Engine, userHandler *userhandler.UserHandler) {
	r.POST("/login", userHandler.Login)
}

func RegisterProductRoutes(r *gin.Engine, productHandler *hproduct.Handler, userUsecase *userusecase.UserUsecase) {
	group := r.Group("/products")
	group.Use(middleware.JWT(userUsecase))
	group.GET("/", productHandler.GetAll)
	group.GET("/:id", productHandler.GetByID)
	group.POST("/", productHandler.Create)
	group.PATCH("/:id", productHandler.Update)
	group.PATCH("/:id/deactivate", productHandler.Delete)
}

func RegisterSupplierRoutes(r *gin.Engine, supplierHandler *hsupplier.Handler, userUsecase *userusecase.UserUsecase) {
	group := r.Group("/suppliers")
	group.Use(middleware.JWT(userUsecase))
	group.GET("/", supplierHandler.GetAll)
	group.GET("/:id", supplierHandler.GetByID)
	group.POST("/", supplierHandler.Create)
	group.PATCH("/:id", supplierHandler.Update)
	group.PATCH("/:id/deactivate", supplierHandler.Delete)
}

func RegisterWarehouseRoutes(r *gin.Engine, warehouseHandler *hwarehouse.Handler, userUsecase *userusecase.UserUsecase) {
	group := r.Group("/warehouses")
	group.Use(middleware.JWT(userUsecase))
	group.GET("/", warehouseHandler.GetAll)
	group.GET("/:id", warehouseHandler.GetByID)
	group.POST("/", warehouseHandler.Create)
	group.PATCH("/:id", warehouseHandler.Update)
	group.PATCH("/:id/deactivate", warehouseHandler.Delete)
}

func RegisterInventoryRoutes(r *gin.Engine, inventoryHandler *hinventory.Handler, userUsecase *userusecase.UserUsecase) {
	group := r.Group("/inventories")
	group.Use(middleware.JWT(userUsecase))
	group.GET("/", inventoryHandler.GetAll)
	group.GET("/:id", inventoryHandler.GetByID)
	group.GET("/product/:product_id", inventoryHandler.GetByProduct)
	group.GET("/warehouse/:warehouse_id", inventoryHandler.GetByWarehouse)
	group.POST("/", inventoryHandler.Create)
	group.PATCH("/:id", inventoryHandler.Update)
	group.DELETE("/:id", inventoryHandler.Delete)
}

func RegisterInventoryMovementRoutes(r *gin.Engine, inventoryMovementHandler *inventorymovement.Handler, userUsecase *userusecase.UserUsecase) {
	group := r.Group("/inventory-movements")
	group.Use(middleware.JWT(userUsecase))
	group.GET("/", inventoryMovementHandler.GetAll)
}

func RegisterPurchaseRequestRoutes(r *gin.Engine, purchaseRequestHandler *hpurchaserequest.Handler, userUsecase *userusecase.UserUsecase) {
	group := r.Group("/purchase-requests")
	group.Use(middleware.JWT(userUsecase))
	group.GET("/", purchaseRequestHandler.GetAll)
	group.GET("/:id", purchaseRequestHandler.GetByID)
	group.POST("/", purchaseRequestHandler.Create)
	group.PATCH("/:id/approve", purchaseRequestHandler.Approve)
	group.PATCH("/:id/reject", purchaseRequestHandler.Reject)
	group.PATCH("/:id", purchaseRequestHandler.Update)
}

func RegisterPurchaseOrderRoutes(r *gin.Engine, purchaseOrderHandler *hpurchaseorder.Handler, userUsecase *userusecase.UserUsecase) {
	group := r.Group("/purchase-orders")
	group.Use(middleware.JWT(userUsecase))
	group.GET("/", purchaseOrderHandler.GetAll)
	group.GET("/purchase-request/:purchase_request_id", purchaseOrderHandler.GetByPurchaseRequestID)
	group.GET("/:id", purchaseOrderHandler.GetByID)
	group.PATCH("/:id/status", purchaseOrderHandler.UpdateStatus)
	r.POST("/purchase-requests/:purchase_request_id/purchase-order", middleware.JWT(userUsecase), purchaseOrderHandler.CreateFromPurchaseRequest)
}

func RegisterGoodsReceiptRoutes(r *gin.Engine, goodsReceiptHandler *hgoodsreceipt.Handler, userUsecase *userusecase.UserUsecase) {
	purchaseOrderGroup := r.Group("/purchase-orders")
	purchaseOrderGroup.Use(middleware.JWT(userUsecase))
	purchaseOrderGroup.POST("/:purchase_order_id/goods-receipts", goodsReceiptHandler.Create)

	goodsReceiptGroup := r.Group("/goods-receipts")
	goodsReceiptGroup.Use(middleware.JWT(userUsecase))
	goodsReceiptGroup.GET("/", goodsReceiptHandler.GetAll)
	goodsReceiptGroup.GET("/:id", goodsReceiptHandler.GetByID)
}
