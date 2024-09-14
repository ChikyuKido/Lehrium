package models

import (
	"golang.org/x/crypto/bcrypt"
)

func (user *User) HashPassword(password string) error {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
  if err != nil {
    return err
  }
  user.Password = string(bytes)
  return nil
}

func (user *User) CheckPassword(providedPassword string) error {
  err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
  if err != nil {
    return err
  }
  return nil
}
