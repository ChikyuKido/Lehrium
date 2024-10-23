package controller

import (
	"lehrium-backend/internal/auth"
	"lehrium-backend/internal/repo"
	"lehrium-backend/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var request struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		UntisName   string `json:"untisName"`
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if repo.DoesUserByEmailExists(request.Email) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		context.Abort()
		return
	}

	hashedPassword, err := util.HashPassword(request.Password)
    if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

    request.Password = hashedPassword

	repo.CreateNewUser(request.Email, request.Password, request.UntisName)

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

	if err := util.CheckPassword(user, request.Password); err != nil {
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
