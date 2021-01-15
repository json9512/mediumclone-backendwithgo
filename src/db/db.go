package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/now"
	_ "github.com/lib/pq"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal("godotenv failed to load variable", err)
	}
}

// ConnectDB ...
// Returns the AwS RDS postgresql database
func ConnectDB() (*gorm.DB, string, error) {
	// Load configuration from util
	DBHost, err := config.LoadConfig("DB_HOST")
	CheckErr(err)
	DBPort, err := config.LoadConfig("DB_PORT")
	CheckErr(err)
	DBName, err := config.LoadConfig("DB_NAME")
	CheckErr(err)
	DBUsername, err := config.LoadConfig("DB_USERNAME")
	CheckErr(err)
	DBPassword, err := config.LoadConfig("DB_PASSWORD")
	CheckErr(err)

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
		log.Println(rdsConnectionString) // TODO: delete; show connection variables for testing
		return nil, "Connection to AWS RDS DB Failed", err
	}

	return db, "Connection to AWS RDS DB Successful", nil
}
