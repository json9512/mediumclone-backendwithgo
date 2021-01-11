package db

import (
	"fmt"

	util "json9512/mediumclone-go/util"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// ConnectDB ...
// Returns the AwS RDS postgresql database
func ConnectDB() (*gorm.DB, string, error) {
	// Load configuration from util
	config, err := util.LoadConfig(".")
	if err != nil {
		fmt.Println("[ERROR] Failed to config")
	}

	// Construct rdsConnectionString with Database configuration
	rdsConnectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s",
		config.Database.DBHost,
		config.Database.DBPort,
		config.Database.DBName,
		config.Database.DBUsername,
		config.Database.DBPassword,
	)

	db, err := gorm.Open("postgres", rdsConnectionString)

	if err != nil {
		return nil, "Connection to AWS RDS DB Failed", err
	}

	return db, "Connection to AWS RDS DB Successful", nil
}
