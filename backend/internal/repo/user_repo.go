package repo

import (
	"fmt"
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
)

func DoesUserByEmailExists(email string) bool {
    var count int64 = 0
	if err := database.New().Instance().Model(models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		fmt.Println(err)
        return true
	}
	return count != 0
}

func CreateNewUser(email, password, untisName string) error {
	user := models.User{
		Email:     email,
		Password:  password,
		UntisName: untisName,
	}
	return database.New().Instance().Create(&user).Error
}

func GetAuthEntry(userid uint) (models.Verification, error){
    var verification models.Verification
	if err := database.New().Instance().Where("user_id = ?", userid).First(&verification).Error; err != nil {
        return verification, err
    }
    return verification, nil
}

func GetUser(email string) (models.User, error) {
	var user models.User
	if err := database.New().Instance().Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

