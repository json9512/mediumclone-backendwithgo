package users

import (
	"github.com/jinzhu/gorm"
)

// User depicts users table in the database
type User struct {
	gorm.Model
	Email    string `gorm:"column:email"`
	Username string `gorm:"column:username"`
}

// CreateTestSample creates a sample user in the database
func CreateTestSample(db *gorm.DB) {
	email := "test@test.com"
	username := "test-user"

	user := User{
		Email:    email,
		Username: username,
	}

	db.Create(&user)
}
