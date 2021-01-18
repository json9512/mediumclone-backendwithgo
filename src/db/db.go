package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// dependencies for above package
	_ "github.com/jinzhu/now"
	_ "github.com/lib/pq"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/logger"
)

// Database ...
// holds reference to gorm.DB object
type Database struct {
	*gorm.DB
}

// DB ...
// global DB variable for export
var DB *gorm.DB

// Init ...
// Returns the AwS RDS postgresql database
func Init() *gorm.DB {
	log := logger.InitLogger()

	// Load configuration from util
	DBHost := config.LoadConfig("DB_HOST")
	DBPort := config.LoadConfig("DB_PORT")
	DBName := config.LoadConfig("DB_NAME")
	DBUsername := config.LoadConfig("DB_USERNAME")
	DBPassword := config.LoadConfig("DB_PASSWORD")

	// Construct rdsConnectionString with Database configuration
	rdsConnectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s",
		DBHost,
		DBPort,
		DBName,
		DBUsername,
		DBPassword,
	)

	db, err := gorm.Open("postgres", rdsConnectionString)

	if err != nil {
		log.Fatal("Connection to AWS RDS DB failed", err)
	} else {
		log.Info("DB connection successful")
	}

	DB = db
	return DB
}

// TestDBInit ...
// returns a DB instance for testing
func TestDBInit() *gorm.DB {
	log := logger.InitLogger()

	dsn := "host=localhost user=postgres password=postgres dbname=mediumclone port=5432 sslmode=disable"
	testDB, err := gorm.Open("postgres", dsn)

	if err != nil {
		log.Fatal("Connection to Test DB failed", err)
	}

	testDB.LogMode(true)
	return testDB
}

// GetDB ...
// Use this function to get connected and serve a pool
func GetDB() *gorm.DB {
	return DB
}
