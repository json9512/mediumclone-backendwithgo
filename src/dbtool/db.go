package dbtool

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// dependencies for above package
	_ "github.com/jinzhu/now"
	_ "github.com/lib/pq"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
)

// Pool manages the gorm.DB struct
type Pool struct {
	*gorm.DB
}

// Init returns the db when connected gracefully
func Init() *Pool {
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

	return &Pool{
		db,
	}
}

// Migrate creates necessary tables in db
func Migrate(db *Pool) {
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&User{})
}

// Query finds the given record in db
func (p *Pool) Query(obj interface{}, condition map[string]interface{}) error {
	query := p.Where(condition).Find(obj)
	return checkErr(query)
}

// Insert creates a new record in db
func (p *Pool) Insert(obj interface{}) error {
	query := p.Create(obj)
	return checkErr(query)
}

// Update updates the record in db
func (p *Pool) Update(user interface{}) error {
	query := p.Model(user).Updates(user)
	return checkErr(query)
}

// Delete hard deletes the record in db
func (p *Pool) Delete(obj interface{}, condition map[string]interface{}) error {
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
