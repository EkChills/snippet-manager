package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type Snippet struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"not null" json:"title"`
    Language  string    `gorm:"not null" json:"language"`
    Content   string    `gorm:"not null" json:"content"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    CreatedAt time.Time
}

