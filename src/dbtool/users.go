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

// GetUserByID gets the user from the database with the given ID
func (db *DB) GetUserByID(id interface{}) (*User, error) {
	var user User
	query := db.First(&user, "id = ?", id)

	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail gets the user from the database with the given email
func (db *DB) GetUserByEmail(email string) (*User, error) {
	var user User
	query := db.First(&user, "email = ?", email)

	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a user with the given credentials in the database
func (db *DB) CreateUser(email, pwd string) (*User, error) {
	user := User{Email: email, Password: pwd}
	query := db.Create(&user)

	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckIfUserExists checks if the user with the given ID exists
func (db *DB) CheckIfUserExists(id interface{}) bool {
	_, err := db.GetUserByID(id)
	if err != nil {
		return false
	}
	return true
}

// UpdateUser updates the user with the provided data
func (db *DB) UpdateUser(newData interface{}) (*User, error) {
	var updatedUser User
	query := db.Model(&updatedUser).Updates(newData)
	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

// DeleteUserByID deletes the user with the given ID in DB
func (db *DB) DeleteUserByID(id interface{}) (*User, error) {
	user, err := db.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	query := db.Unscoped().Delete(user)
	if err = checkErr(query); err != nil {
		return nil, err
	}
	return user, nil
}
