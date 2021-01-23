package users

import (
	"github.com/jinzhu/gorm"
)

// User depicts users table in the database
type User struct {
	gorm.Model
	Email        string `gorm:"column:email;not null;unique"`
	Password     string `gorm:"column:password;not null"`
	AccessToken  string `gorm:"column:access_token"`
	RefreshToken string `gorm:"column:refresh_token"`
}

// CreateTestSample creates a sample user in the database
func CreateTestSample(db *gorm.DB) {
	user := User{
		Email:        "test@test.com",
		Password:     "test-password",
		AccessToken:  "",
		RefreshToken: "",
	}

	db.Create(&user)
}
