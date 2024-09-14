package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
    "github.com/lib/pq"
)

type User struct {
    gorm.Model
    ID              uint    `gorm:"primaryKey"`
    Email           string  `json:"email" gorm:"size:255;unique"`
    Password        string  `json:"password" gorm:"size:255"`
    UntisName       string  `json:"untisName" gorm:"size:100;unique"`
    Roles           pq.StringArray `gorm:"type:varchar(50)[]"`
    TeacherIDs      pq.Int64Array `gorm:"type:integer[]"`
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
