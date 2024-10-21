package controller

import (
	"fmt"
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	var user models.User
	var verification models.Verification

	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no UUID provided"})
	}

	usersRecord := database.New().Instance().Where("uuid = ?", uuid).First(&user)
	if usersRecord.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": usersRecord.Error.Error()})
		c.Abort()
		return
	}

	verificationRecord := database.New().Instance().Where("user_id = ?", user.ID).First(&verification)
	if verificationRecord.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": verificationRecord.Error})
		c.Abort()
		return
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", verification.ExpDate)
	if err != nil {
		fmt.Println(verification.ExpDate)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check if link expired"})
		return
	}

	if parsedTime.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification Link expired"})
		c.Abort()
		return
	}

	if user.Verified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already verified"})
		c.Abort()
		return
	}
	database.New().Instance().Model(&models.User{}).Where("email = ?", user.Email).Update("isVerified", true)
	return
}

func SendVerificationEmail(uuid string, email string) {
	from := os.Getenv("GMAIL_USERNAME")
	pass := os.Getenv("GMAIL_APPPASSWORD")
	to := email

	msg := "From: Lehrium Verification" + "\n" +
		"To: " + to + "\n" +
		"Subject: Lehrium Account verification\n\n" +
		"please verify your account via this link: \n" +
		"https://lehrium.elekius.at/auth/verifyEmail?uuid=" + uuid

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Println("Successfully sended to " + to)
}

// updates ExpirationDate with current Time + 5 Minutes
func updateVerificationEmailDate(uuid string) {
	var user models.User

	usersRecord := database.New().Instance().Where("uuid = ?", uuid).First(&user)
	if usersRecord.Error != nil {
		return
	}
	database.New().Instance().Model(&models.Verification{}).Where("user_id = ?", user.ID).Update("exp_date", time.Now().Add(5*time.Minute).String())
}
