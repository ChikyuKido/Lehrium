package util

import (
	"lehrium-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (hashedPassword string, err error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
  if err != nil {
    return "", err
  }
  return string(bytes), nil
}

func CheckPassword(user models.User, providedPassword string) error {
  err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
  if err != nil {
    return err
  }
  return nil
}
