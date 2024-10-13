package middlewares

import (
  "lehrium-backend/internal/auth"
  "github.com/gin-gonic/gin"
  "strings"
)

func Auth() gin.HandlerFunc{
    return func(context *gin.Context) {
        tokenString := context.GetHeader("Authorization")
        if tokenString == "" {
            context.JSON(401, gin.H{"error": "request does not contain an access token"})
            context.Abort()
            return
        }

        // Check if the Authorization header starts with "Bearer"
        if !strings.HasPrefix(tokenString, "Bearer ") {
            context.JSON(401, gin.H{"error": "invalid authorization header format"})
            context.Abort()
            return
        }

        // Extract the actual token by removing the "Bearer " prefix
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        err := auth.ValidateToken(tokenString)
        if err != nil {
            context.JSON(401, gin.H{"error": err.Error()})
            context.Abort()
            return
        }

        context.Next()  
    }
}
