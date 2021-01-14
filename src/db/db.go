package db

import (
	"fmt"
	"log"

	"github.com/json9512/mediumclone-backendwithgo/src/util"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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
	DBHost, err := util.LoadConfig("DB_HOST")
	CheckErr(err)
	DBPort, err := util.LoadConfig("DB_PORT")
	CheckErr(err)
	DBName, err := util.LoadConfig("DB_NAME")
	CheckErr(err)
	DBUsername, err := util.LoadConfig("DB_USERNAME")
	CheckErr(err)
	DBPassword, err := util.LoadConfig("DB_PASSWORD")
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
		return nil, "Connection to AWS RDS DB Failed", err
	}

	return db, "Connection to AWS RDS DB Successful", nil
}
