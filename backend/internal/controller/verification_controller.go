package controller

import (
	"fmt"
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
	"lehrium-backend/internal/repo"
	"net/http"
	"time"
    "log"
    "net/smtp"
    "os"

	"github.com/gin-gonic/gin"
)

func SendVerificationEmail(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        fmt.Println("failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

    var verification models.Verification
    verification, err := repo.GetAuthEntry(user.ID)
    if err != nil {
        log.Panicln("failed to retrieve authentication record")
    }
    
    from := os.Getenv("GMAIL_USERNAME")
	pass := os.Getenv("GMAIL_APPPASSWORD")
    baseurl := os.Getenv("BASE_URL")
    smtpServerUrl := os.Getenv("SMTP_SERVER")
    smtpServerPort := os.Getenv("SMTP_PORT")
	to := user.Email

	msg := "From: Lehrium Verification" + "\n" +
		"To: " + to + "\n" +
		"Subject: Lehrium Account verification\n\n" +
		"please verify your account via this link: \n" +
		"https://" + baseurl + "/auth/verifyEmail?uuid=" + verification.UUID

        err = smtp.SendMail(fmt.Sprintf("%s:%s",smtpServerUrl, smtpServerPort),
		smtp.PlainAuth("", from, pass, smtpServerUrl),
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
