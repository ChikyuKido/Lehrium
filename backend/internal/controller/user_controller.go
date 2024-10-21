package controller

import (
	"lehrium-backend/internal/auth"
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
	"lehrium-backend/internal/repo"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

func RegisterUser(context *gin.Context) {
	var user models.User
	var userAuth models.Verification

	userAuth.UUID = uuid.NewString()
	user.Verified = false

	userAuth.UserID = user.ID
	userAuth.ExpDate = time.Now().Add(time.Minute * 5).String()
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// check auf Email Domain
	/*
		match, err := repo.CheckForEmailDomain(user.Email)
		if err != nil || !match {
			context.JSON(http.StatusBadRequest, gin.H{"error": err})
			context.Abort()
			return
		}
	*/
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Create the new user usersRecord in the database
	usersRecord := database.New().Instance().Create(&user)
	if usersRecord.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": usersRecord.Error.Error()})
		context.Abort()
		return
	}
	authenticationRecord := database.New().Instance().Create(&userAuth)
	if authenticationRecord.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": authenticationRecord.Error.Error()})
	}

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
