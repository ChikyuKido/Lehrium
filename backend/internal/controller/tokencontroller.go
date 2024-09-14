package controllers

import (
	"lehrium-backend/internal/auth"
	//"backend/internal/database"
	"lehrium-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)
type TokenRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
func GenerateToken(c *gin.Context) {
    var request TokenRequest
    var user models.User
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        c.Abort()
        return
    }

    // need to get user from db.
    user = models.User{Email: "test@test.com", Password: "$2a$07$QU.NPhhUet6shMMrW0cYEOYsXCHrmU5iCrysowxadRuTOLjoDtRzC"/* hashed password go here */ }
    credentialError := user.CheckPassword(request.Password)
    if credentialError != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        c.Abort()
        return
    }
    tokenString, err:= auth.GenerateJWT(user.Email, user.UntisName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        c.Abort()
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
