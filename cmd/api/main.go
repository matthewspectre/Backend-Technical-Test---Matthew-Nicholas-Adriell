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

	"be_evindo/internal/handler"
	producthandler "be_evindo/internal/handler/product"
	postgresrepository "be_evindo/internal/postgres"
	productrepository "be_evindo/internal/postgres/product"
	"be_evindo/internal/router"
	"be_evindo/internal/usecase"
	productusecase "be_evindo/internal/usecase/product"
)

func main() {
	dsn := postgresDSN()

	database, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDatabase, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDatabase.Close()

	userRepository := postgresrepository.NewUserRepository(sqlDatabase)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	productRepository := productrepository.NewRepository(database)
	productUsecase := productusecase.NewUsecase(productRepository)
	productHandler := producthandler.NewHandler(productUsecase)

	engine := gin.Default()
	engine.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusNoContent)
	})
	legacyHandler := router.New(userHandler)
	engine.GET("/users/:id", gin.WrapH(legacyHandler))
	engine.POST("/users", gin.WrapH(legacyHandler))
	router.RegisterProductRoutes(engine, productHandler)

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
