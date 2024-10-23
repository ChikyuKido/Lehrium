package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Teacher struct {
    gorm.Model
    Name      string `gorm:"size:50"`
    MiddleName string `gorm:"size:50"`
    LastName  string `gorm:"size:50"`
    ShortName string `gorm:"size:3;unique"`
    Comments  []Comment
    Ratings   []Rating
}

// Comment represents the "Comments" table
type Comment struct {
    gorm.Model
    ID          uint      `gorm:"primaryKey"`
    Content     string    `gorm:"type:text"`
    TeacherID   uint      `gorm:"index"`
    CreationDate time.Time `gorm:"type:date"`
    Teacher     Teacher   `gorm:"foreignKey:TeacherID"`
}

// Rating represents the "Ratings" table
type Rating struct {
    gorm.Model
    TeacherID      uint      `gorm:"index"`
    CreationDate   time.Time `gorm:"type:date"`
    TeachingSkills int       `gorm:"check:teaching_skills >= 1 AND teaching_skills <= 5"`
    Kindness       int       `gorm:"check:kindness >= 1 AND kindness <= 5"`
    Engagement     int       `gorm:"check:engagement >= 1 AND engagement <= 5"`
    Organization   int       `gorm:"check:organization >= 1 AND organization <= 5"`
    Teacher        Teacher   `gorm:"foreignKey:TeacherID"`
}

type User struct {
    gorm.Model
    Email           string  `json:"email" gorm:"size:255;unique"`
    Password        string  `json:"password" gorm:"size:255"`
    UntisName       string  `json:"untisName" gorm:"size:100;unique"`
    Roles           pq.StringArray `gorm:"type:varchar(50)[]"`
    TeacherIDs      pq.Int64Array `gorm:"type:integer[]"`
    Verified        bool    `json:"isVerified" gorm:"type:boolean"`
    RatedTeachers   []bool  `json:"ratedTeacehrs" gorm:"type:bool[]"`
}

type Verification struct {
    gorm.Model
    UserID  uint
    UUID    string
    ExpDate int64
}
