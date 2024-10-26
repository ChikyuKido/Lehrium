package controller

import (
	"lehrium-backend/internal/auth"
	"lehrium-backend/internal/repo"
	"lehrium-backend/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var request struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		UntisName   string `json:"untisName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body json. Contact the frontend dev"})
		c.Abort()
		return
	}

	if repo.DoesUserByEmailExists(request.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		c.Abort()
		return
	}

	hashedPassword, err := util.HashPassword(request.Password)
    if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash user password"})
		c.Abort()
		return
	}

    request.Password = hashedPassword

    err = repo.CreateNewUser(request.Email, request.Password, request.UntisName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account due to internal error. Try again later"})
    }

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully created"})
}

func LoginUser(c *gin.Context) {
	var request struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberme"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body json. Contact the frontend dev"})
		c.Abort()
		return
	}

	user, err := repo.GetUser(request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		c.Abort()
		return
	}

	if err := util.CheckPassword(user, request.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		c.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.Email, request.RememberMe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log in"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
