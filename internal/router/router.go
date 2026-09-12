package router

import (
	"net/http"

	"be_evindo/internal/handler"
	hproduct "be_evindo/internal/handler/product"

	"github.com/gin-gonic/gin"
)

func New(userHandler *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("POST /users", userHandler.Create)
	return mux
}

func RegisterProductRoutes(r *gin.Engine, productHandler *hproduct.Handler) {
	group := r.Group("/products")
	group.POST("/", productHandler.Create)
	group.GET("/", productHandler.GetAll)
	group.GET("/:id", productHandler.GetByID)
	group.PATCH("/:id", productHandler.Update)
	group.PATCH("/:id/deactivate", productHandler.Delete)
}
