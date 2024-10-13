package controller

import (
	"lehrium-backend/internal/auth"
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
	"net/http"
	"time"

	//    "regexp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type TokenRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberme"`
}

func RegisterUser(context *gin.Context) {
	var user models.User
	var userAuth models.Verification

	user.UUID = uuid.NewString()
	user.Verified = false

	userAuth.UUID = user.UUID
	userAuth.UserID = user.ID
	userAuth.ExpDate = time.Now().Add(time.Minute * 5).String()
	// Bind the incoming JSON data to the user struct
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// check auf Email Domain
	/*
	   email_match, err := regexp.MatchString(`@spengergasse\.at$`, user.Email)
	   if !email_match {
	       context.JSON(http.StatusBadRequest, gin.H{"error": "this email is not part of the school mail network"})
	       context.Abort()
	       return
	   }
	   if err != nil {
	       context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	       context.Abort()
	       return
	   }
	*/
	// Hash the password before saving it to the database
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Create the new user usersRecord in the database
	usersRecord := database.New().Instance().Create(&user) // Get the instance from the database package and save the user
	if usersRecord.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": usersRecord.Error.Error()})
		context.Abort()
		return
	}
	authenticationRecord := database.New().Instance().Create(&userAuth)
	if authenticationRecord.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": authenticationRecord.Error.Error()})
	}
	sendFirstVerificationEmail(user.UUID, user.Email)

	// If the user is successfully created, return a success response
	context.JSON(http.StatusCreated, gin.H{"message": "Successfully created"})
}

func LoginUser(c *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	usersRecord := database.New().Instance().Where("email = ?", request.Email).First(&user)
	if usersRecord.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": usersRecord.Error.Error()})
		c.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
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
