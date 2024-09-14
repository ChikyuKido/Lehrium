package models

import (
    "time"
)

type Teacher struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"size:50"`
    MiddleName string `gorm:"size:50"`
    LastName  string `gorm:"size:50"`
    ShortName string `gorm:"size:3;unique"`
    Comments  []Comment
    Ratings   []Rating
}

// Comment represents the "Comments" table
type Comment struct {
    ID          uint      `gorm:"primaryKey"`
    Content     string    `gorm:"type:text"`
    TeacherID   uint      `gorm:"index"`
    CreationDate time.Time `gorm:"type:date"`
    Teacher     Teacher   `gorm:"foreignKey:TeacherID"`
}

// Rating represents the "Ratings" table
type Rating struct {
    ID             uint      `gorm:"primaryKey"`
    TeacherID      uint      `gorm:"index"`
    CreationDate   time.Time `gorm:"type:date"`
    TeachingSkills int       `gorm:"check:teachingSkills >= 1 AND teachingSkills <= 10"`
    Kindness       int       `gorm:"check:kindness >= 1 AND kindness <= 10"`
    Engagement     int       `gorm:"check:engagement >= 1 AND engagement <= 10"`
    Organization   int       `gorm:"check:organization >= 1 AND organization <= 10"`
    Teacher        Teacher   `gorm:"foreignKey:TeacherID"`
}
