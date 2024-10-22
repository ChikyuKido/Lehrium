package controller

import (
	"fmt"
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
	"lehrium-backend/internal/repo"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SendVerificationEmail(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "error casting email to string"})
		c.Abort()
		return
	}

	user, err := repo.GetUser(emailStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	uuid := uuid.NewString()
	repo.CreateNewAuthenticationRecord(user.ID, uuid)

	from := os.Getenv("GMAIL_USERNAME")
	pass := os.Getenv("GMAIL_APPPASSWORD")
	baseurl := os.Getenv("BASE_URL")
	//smtpServerUrl := os.Getenv("SMTP_SERVER")
	//smtpServerPort := os.Getenv("SMTP_PORT")
	to := user.Email

	msg := "From: Lehrium Verification" + "\n" +
		"To: " + to + "\n" +
		"Subject: Lehrium Account verification\n\n" +
		"please verify your account via this link: \n" +
		"https://" + baseurl + "/api/v1/auth/verifyEmail?uuid=" + uuid

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Println("Successfully sended to " + to)
}

func VerifyEmail(c *gin.Context) {
	var user models.User
	var verification models.Verification

	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no UUID provided"})
	}

	verificationRecord := database.New().Instance().Where("uuid = ?", uuid).First(&verification)
	if verificationRecord.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": verificationRecord.Error})
		c.Abort()
		return
	}

	usersRecord := database.New().Instance().Where("id = ?", verification.UserID).First(&user)
	if usersRecord.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": usersRecord.Error.Error()})
		c.Abort()
		return
	}

	fmt.Println(verification.ExpDate)
	// bin am verzwifeln
	layout := ""
	parsedTime, err := time.Parse(layout, verification.ExpDate[:40])
	if err != nil {
		fmt.Println(err)
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
