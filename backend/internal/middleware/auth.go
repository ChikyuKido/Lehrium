package middlewares

import (
	"lehrium-backend/internal/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			context.JSON(401, gin.H{"error": "invalid authorization header format"})
			context.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		context.Set("email", claims.Email)

		context.Next()
	}
}
