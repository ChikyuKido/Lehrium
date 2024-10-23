package repo

import (
	"errors"
	"fmt"
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
	"log"
	"time"
)

func CreateNewAuthenticationRecord(userid uint, uuid string) {
	var verification = models.Verification{
		UserID:  userid,
		UUID:    uuid,
		ExpDate: time.Now().Add(time.Minute * 5).Unix(),
	}

	if err := database.New().Instance().Create(&verification).Error; err != nil {
		log.Panicln("failed to create authentication")
	}
}

func VerifyUser(verification models.Verification, user models.User) error{
	if user.Verified {
        err := errors.New("user is already verified")
		return err
	}
	database.New().Instance().Model(&models.User{}).Where("email = ?", user.Email).Update("verified", true)
    return nil
}

func CheckIfAuthenticationRecordExists(userid uint) bool{
    var count int64 = 0
    if err := database.New().Instance().Where("user_id = ?", userid).Count(&count).Error; err != nil {
        fmt.Println(err)
        return true
    }
	return count != 0
}

func CheckIfAuthenticationDateExpired(userid uint) bool{
    var verification models.Verification

    verificationRecord := database.New().Instance().Where("user_id = ?", userid).First(&verification)
    if verificationRecord.Error != nil {
        return true
    }

    expDate := time.Unix(verification.ExpDate, 0)
    if time.Now().After(expDate){
        return true
    }

    return false
}

func DeleteAuthenticationRecord(userid uint){
    database.New().Instance().Delete(models.Verification{UserID: userid})
}
