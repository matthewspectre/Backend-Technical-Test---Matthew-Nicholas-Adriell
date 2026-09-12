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

	producthandler "be_evindo/internal/handler/product"
	userhandler "be_evindo/internal/handler/user"
	productrepository "be_evindo/internal/postgres/product"
	userrepository "be_evindo/internal/postgres/user"
	"be_evindo/internal/router"
	productusecase "be_evindo/internal/usecase/product"
	userusecase "be_evindo/internal/usecase/user"
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
	userRepository := userrepository.NewRepository(database)
	userUsecase := userusecase.NewUserUsecase(userRepository, envOrDefault("JWT_SECRET", "be-evindo-development-secret"))
	userHandler := userhandler.NewUserHandler(userUsecase)

	engine := gin.Default()
	engine.GET("/health", func(context *gin.Context) {
		context.Status(http.StatusNoContent)
	})
	router.RegisterAuthRoutes(engine, userHandler)
	router.RegisterLoginRoute(engine, userHandler)
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
