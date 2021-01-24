package users

import (
	"time"

	"gorm.io/gorm"
)

// User depicts users table in the database
type User struct {
	ID           uint   `gorm:"primary_key"`
	Email        string `gorm:"column:email;not null;unique"`
	Password     string `gorm:"column:password;not null"`
	AccessToken  string `gorm:"column:access_token"`
	RefreshToken string `gorm:"column:refresh_token"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
