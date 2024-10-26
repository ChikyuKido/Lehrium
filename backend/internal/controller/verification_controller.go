package controller

import (
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
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

    user, ok := userInterface.(models.User)
    if !ok {
        c.JSON(401, gin.H{"error": "Internal Error"})
        c.Abort()
        return
    }

	user, err := repo.GetUser(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User with provided email does not exist"})
        c.Abort()
        return
	}

    if repo.CheckIfAuthenticationRecordExists(user.ID){
        if repo.CheckIfAuthenticationDateExpired(user.ID){
            repo.DeleteAuthenticationRecord(user.ID)
        }else {   
            c.JSON(http.StatusInternalServerError, gin.H{"error": "An email verification is already pending"})
            c.Abort()
            return
        }
    }
    uuid := uuid.NewString()
    repo.CreateNewAuthenticationRecord(user.ID, uuid)



	from := os.Getenv("GMAIL_USERNAME")
	pass := os.Getenv("GMAIL_APPPASSWORD")
	baseurl := os.Getenv("FRONTEND_URL")
	//smtpServerUrl := os.Getenv("SMTP_SERVER")
	//smtpServerPort := os.Getenv("SMTP_PORT")
	to := user.Email

	msg := "From: Lehrium Verification" + "\n" +
		"To: " + to + "\n" +
		"Subject: Lehrium Account verification\n\n" +
		"please verify your account via this link: \n" +
		baseurl + "/auth/verifyEmail?uuid=" + uuid

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
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no UUID provided"})
	}


    verification, err := repo.GetVerificationByUUID(uuid)
    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get users verification info"})
		c.Abort()
		return
    }

    user, err := repo.GetUserById(verification.ID)
    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get user"})
        c.Abort()
        return
    }
    
    expDate := time.Unix(verification.ExpDate, 0)
    if time.Now().After(expDate){
        c.JSON(400, gin.H{"error": "verification link expired"})
        c.Abort()
        return
    }

    err = repo.VerifyUser(verification, user)
    if err != nil {
        c.JSON(400, gin.H{"error": "failed to verify user"})
        c.Abort()
        return
    }
}
