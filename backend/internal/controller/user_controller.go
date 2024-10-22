package controller

import (
	"lehrium-backend/internal/auth"
	"lehrium-backend/internal/models"
	"lehrium-backend/internal/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if repo.DoesUserByEmailExists(user.Email) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		context.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	repo.CreateNewUser(user.Email, user.Password, user.UntisName)

	context.JSON(http.StatusCreated, gin.H{"message": "Successfully created"})
}

func LoginUser(c *gin.Context) {
	var request struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberme"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user, err := repo.GetUser(request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		c.Abort()
		return
	}

	if err := user.CheckPassword(request.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.Email, request.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
