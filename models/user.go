package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   string `json:"user_id" gorm:"type:varchar(20);uniqueIndex"` // Sonyflake ID
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Phone    string `json:"phone"`
	Bio      string `json:"bio"`
}
