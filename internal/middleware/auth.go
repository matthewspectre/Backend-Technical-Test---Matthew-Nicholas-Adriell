package middleware

import (
	"net/http"
	"strings"

	userusecase "be_evindo/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

func JWT(userUsecase *userusecase.UserUsecase) gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization token is required"})
			return
		}

		if _, err := userUsecase.ParseToken(strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))); err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		context.Next()
	}
}
