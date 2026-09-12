package router

import (
	hproduct "be_evindo/internal/handler/product"
	hsupplier "be_evindo/internal/handler/supplier"
	userhandler "be_evindo/internal/handler/user"
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
