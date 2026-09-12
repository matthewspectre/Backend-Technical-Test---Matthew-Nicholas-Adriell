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

		claims, err := userUsecase.ParseToken(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		userID, ok := claims["user_id"].(float64)
		if !ok || userID <= 0 {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user identity"})
			return
		}
		context.Set("user_id", int64(userID))
		context.Set("user_role", claims["role"])
		context.Next()
	}
}
