package models

import (
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID          int
    Email       string  `json:"email"`
    Password    string  `json:"password"`
    UntisName   string  `json:"untisName"`
}

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
