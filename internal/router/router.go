package router

import (
	hproduct "be_evindo/internal/handler/product"
	userhandler "be_evindo/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, userHandler *userhandler.UserHandler) {
	group := r.Group("/auth")
	group.POST("/register", userHandler.Register)
}

func RegisterLoginRoute(r *gin.Engine, userHandler *userhandler.UserHandler) {
	r.POST("/login", userHandler.Login)
}

func RegisterProductRoutes(r *gin.Engine, productHandler *hproduct.Handler) {
	group := r.Group("/products")
	group.GET("/", productHandler.GetAll)
	group.GET("/:id", productHandler.GetByID)
	group.POST("/", productHandler.Create)
	group.PATCH("/:id", productHandler.Update)
	group.PATCH("/:id/deactivate", productHandler.Delete)
}
