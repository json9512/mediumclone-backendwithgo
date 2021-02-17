package dbtool

import (
	"time"

	"gorm.io/gorm"
)

// User depicts users table in the database
type User struct {
	ID              uint   `gorm:"primary_key"`
	Email           string `gorm:"column:email;not null;unique"`
	Password        string `gorm:"column:password;not null"`
	TokenExpiryDate int64  `gorm:"column:token_expiry_date"`
	Posts           []Post `gorm:"foreignKey:Author"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
