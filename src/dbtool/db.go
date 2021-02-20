package dbtool

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// dependencies for above package
	_ "github.com/jinzhu/now"
	_ "github.com/lib/pq"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
)

// DB manages the gorm.DB struct
type DB struct {
	*gorm.DB
}

// Init returns the db when connected gracefully
func Init() *DB {
	log := config.InitLogger()
	config := createConfig()

	// Construct configString for database connection
	configString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBUsername,
		config.DBPassword,
	)

	db, err := gorm.Open("postgres", configString)

	if err != nil {
		if config.DBHost != "localhost" {
			log.Fatal("Connection to AWS RDS DB failed", err)
		}
		log.Fatal("Connection to Test DB failed", err)
	} else {
		log.Info("DB connection successful")
	}

	return &DB{
		db,
	}
}

// Migrate creates necessary tables in db
func Migrate(db *DB) {
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&User{})
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

// Query finds the given record in db
func (db *DB) Query(obj interface{}, condition map[string]interface{}) error {
	query := db.Where(condition).Find(obj)
	return checkErr(query)
}

// Insert creates a new record in db
func (db *DB) Insert(obj interface{}) error {
	query := db.Create(obj)
	return checkErr(query)
}

// Update updates the record in db
func (db *DB) Update(obj interface{}) error {
	query := db.Model(obj).Updates(obj)
	return checkErr(query)
}

// Delete hard deletes the record in db
func (db *DB) Delete(obj interface{}, condition map[string]interface{}) error {
	// Soft delete the user
	// query := p.Where(condition).Find(obj).Delete(obj)
	query := db.Unscoped().Delete(obj)
	return checkErr(query)
}

func checkErr(g *gorm.DB) error {
	if g.Error != nil {
		return g.Error
	}
	return nil
}
