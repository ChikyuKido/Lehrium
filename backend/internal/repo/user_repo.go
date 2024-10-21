package repo

import (
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
	"regexp"
)

func DoesUserByEmailExists(email string) bool {
	var count int64
	if err := database.New().Instance().Where("email = ?", email).Count(&count).Error; err != nil {
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

func GetUser(email string) (models.User, error) {
	var user models.User
	if err := database.New().Instance().Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func CheckForEmailDomain(email string) (bool, error) {
	email_match, err := regexp.MatchString(`@spengergasse\.at$`, email)
	if !email_match {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
