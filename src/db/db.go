package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// dependency for above package
	_ "github.com/jinzhu/now"
	_ "github.com/lib/pq"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
)

// ConnectDB ...
// Returns the AwS RDS postgresql database
func ConnectDB() (*gorm.DB, string, error) {
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
		return nil, "Connection to AWS RDS DB Failed", err
	}

	return db, "Connection to AWS RDS DB Successful", nil
}
