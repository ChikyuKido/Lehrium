package repo

import (
	"lehrium-backend/internal/database"
	"lehrium-backend/internal/models"
    "regexp"
)

func DoesUserByEmailExists(email string) bool{
    var count int64;
    if err := database.New().Instance().Where("email = ?", email).Count(&count).Error; err != nil {
        return true;
    } 
    return count != 0;
}

func CreateNewUser(email, password , untisName string) error{

    user := models.User{
        Email: email,
        Password: password,
        UntisName: untisName,
    }

    return database.New().Instance().Create(&user).Error
}

func CheckForEmailDomain(email string) (bool, error){
    if err := email_match, err := regexp.MatchString(`@spengergasse\.at$`, email); err != nil {

    }
	if !email_match {
      return false, nil
    }

}
