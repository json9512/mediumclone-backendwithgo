package users

import (
	"github.com/jinzhu/gorm"
)

// UserModel depicts user_model table in the database
type UserModel struct {
	gorm.Model
	Email    string `gorm:"column:email"`
	Username string `gorm:"column:username"`
}

// CreateTestSample creates a sample user in the database
func CreateTestSample(db *gorm.DB) {
	email := "test@test.com"
	username := "test-user"

	user := UserModel{
		Email:    email,
		Username: username,
	}

	db.Create(&user)
}
