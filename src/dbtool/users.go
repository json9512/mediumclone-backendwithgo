package dbtool

import (
	"time"

	"gorm.io/gorm"
)

// User depicts users table in the database
type User struct {
	ID             uint       `gorm:"primary_key"`
	Email          string     `gorm:"column:email;not null;unique"`
	Password       string     `gorm:"column:password;not null"`
	TokenCreatedAt *time.Time `gorm:"column:token_created_at"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
