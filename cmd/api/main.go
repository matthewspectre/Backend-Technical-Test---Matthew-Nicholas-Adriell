package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	goodsreceipthandler "be_evindo/internal/handler/goods_receipt"
	inventoryhandler "be_evindo/internal/handler/inventory"
	producthandler "be_evindo/internal/handler/product"
	purchaseorderhandler "be_evindo/internal/handler/purchase_order"
	purchaserequesthandler "be_evindo/internal/handler/purchase_request"
	supplierhandler "be_evindo/internal/handler/supplier"
	userhandler "be_evindo/internal/handler/user"
	warehousehandler "be_evindo/internal/handler/warehouse"
	goodsreceiptrepository "be_evindo/internal/postgres/goods_receipt"
	inventoryrepository "be_evindo/internal/postgres/inventory"
	productrepository "be_evindo/internal/postgres/product"
	purchaseorderrepository "be_evindo/internal/postgres/purchase_order"
	purchaserequestrepository "be_evindo/internal/postgres/purchase_request"
	supplierrepository "be_evindo/internal/postgres/supplier"
	userrepository "be_evindo/internal/postgres/user"
	warehouserepository "be_evindo/internal/postgres/warehouse"
	"be_evindo/internal/router"
	goodsreceiptusecase "be_evindo/internal/usecase/goods_receipt"
	inventoryusecase "be_evindo/internal/usecase/inventory"
	productusecase "be_evindo/internal/usecase/product"
	purchaseorderusecase "be_evindo/internal/usecase/purchase_order"
	purchaserequestusecase "be_evindo/internal/usecase/purchase_request"
	supplierusecase "be_evindo/internal/usecase/supplier"
	userusecase "be_evindo/internal/usecase/user"
	warehouseusecase "be_evindo/internal/usecase/warehouse"
)

func main() {
	dsn := postgresDSN()

	database, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	productRepository := productrepository.NewRepository(database)
	productUsecase := productusecase.NewUsecase(productRepository)
	productHandler := producthandler.NewHandler(productUsecase)
	supplierRepository := supplierrepository.NewRepository(database)
	supplierUsecase := supplierusecase.NewUsecase(supplierRepository)
	supplierHandler := supplierhandler.NewHandler(supplierUsecase)
	warehouseRepository := warehouserepository.NewRepository(database)
	warehouseUsecase := warehouseusecase.NewUsecase(warehouseRepository)
	warehouseHandler := warehousehandler.NewHandler(warehouseUsecase)
	inventoryRepository := inventoryrepository.NewRepository(database)
	inventoryUsecase := inventoryusecase.NewUsecase(inventoryRepository, warehouseRepository)
	inventoryHandler := inventoryhandler.NewHandler(inventoryUsecase)
	purchaseRequestRepository := purchaserequestrepository.NewRepository(database)
	purchaseRequestUsecase := purchaserequestusecase.NewUsecase(purchaseRequestRepository, warehouseRepository)
	purchaseRequestHandler := purchaserequesthandler.NewHandler(purchaseRequestUsecase)
	purchaseOrderRepository := purchaseorderrepository.NewRepository(database)
	purchaseOrderUsecase := purchaseorderusecase.NewUsecase(purchaseOrderRepository, purchaseRequestRepository, supplierRepository)
	purchaseOrderHandler := purchaseorderhandler.NewHandler(purchaseOrderUsecase)
	goodsReceiptRepository := goodsreceiptrepository.NewRepository(database)
	goodsReceiptUsecase := goodsreceiptusecase.NewUsecase(goodsReceiptRepository, purchaseOrderRepository)
	goodsReceiptHandler := goodsreceipthandler.NewHandler(goodsReceiptUsecase)
	userRepository := userrepository.NewRepository(database)
	userUsecase := userusecase.NewUserUsecase(userRepository, envOrDefault("JWT_SECRET", "be-evindo-development-secret"))
	userHandler := userhandler.NewUserHandler(userUsecase)

	engine := gin.Default()
	engine.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusNoContent)
	})
	router.RegisterAuthRoutes(engine, userHandler)
	router.RegisterLoginRoute(engine, userHandler)
	router.RegisterProductRoutes(engine, productHandler, userUsecase)
	router.RegisterSupplierRoutes(engine, supplierHandler, userUsecase)
	router.RegisterWarehouseRoutes(engine, warehouseHandler, userUsecase)
	router.RegisterInventoryRoutes(engine, inventoryHandler, userUsecase)
	router.RegisterPurchaseRequestRoutes(engine, purchaseRequestHandler, userUsecase)
	router.RegisterPurchaseOrderRoutes(engine, purchaseOrderHandler, userUsecase)
	router.RegisterGoodsReceiptRoutes(engine, goodsReceiptHandler, userUsecase)

	server := &http.Server{
		Addr:    net.JoinHostPort("", envOrDefault("APP_PORT", "8080")),
		Handler: engine,
	}
	log.Printf("server listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func postgresDSN() string {
	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		return dsn
	}

	databaseURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(envOrDefault("DB_USER", "postgres"), os.Getenv("DB_PASS")),
		Host: net.JoinHostPort(
			envOrDefault("DB_HOST", "localhost"),
			envOrDefault("DB_PORT", "5432"),
		),
		Path: "/" + envOrDefault("DB_NAME", "be_evindo"),
	}
	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()
	return databaseURL.String()
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
