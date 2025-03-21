package models

import (
	"gorm.io/gorm"
)

// User struct for database
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}
