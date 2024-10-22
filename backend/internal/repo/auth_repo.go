package repo

import (
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
	"log"
	"time"
)

func CreateNewAuthenticationRecord(userid uint, uuid string) {
	var verification = models.Verification{
		UserID:  userid,
		UUID:    uuid,
		ExpDate: time.Now().Add(time.Minute * 5).String(),
	}

	if err := database.New().Instance().Create(&verification).Error; err != nil {
		log.Panicln("failed to create authentication")
	}
}
