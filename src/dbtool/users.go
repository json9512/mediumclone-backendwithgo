package dbtool

import (
	"time"

	"gorm.io/gorm"
)

// User depicts users table in the database
type User struct {
	ID             uint   `gorm:"primary_key"`
	Email          string `gorm:"column:email;not null;unique"`
	Password       string `gorm:"column:password;not null"`
	TokenExpiresIn int64  `gorm:"column:token_expires_in"`
	Posts          []Post `gorm:"foreignKey:Author"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
