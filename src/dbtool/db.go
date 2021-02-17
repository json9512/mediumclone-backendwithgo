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

// GetUserByID returns a pointer to the user obj with the given ID and the error
func (p *DB) GetUserByID(id int64) (*User, error) {
	var user User
	query := p.First(&user, "id = ?", id)
	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &user, nil
}

// Query finds the given record in db
func (p *DB) Query(obj interface{}, condition map[string]interface{}) error {
	query := p.Where(condition).Find(obj)
	return checkErr(query)
}

// Insert creates a new record in db
func (p *DB) Insert(obj interface{}) error {
	query := p.Create(obj)
	return checkErr(query)
}

// Update updates the record in db
func (p *DB) Update(obj interface{}) error {
	query := p.Model(obj).Updates(obj)
	return checkErr(query)
}

// Delete hard deletes the record in db
func (p *DB) Delete(obj interface{}, condition map[string]interface{}) error {
	// Soft delete the user
	// query := p.Where(condition).Find(obj).Delete(obj)
	query := p.Unscoped().Delete(obj)
	return checkErr(query)
}

func checkErr(g *gorm.DB) error {
	if g.Error != nil {
		return g.Error
	}
	return nil
}
