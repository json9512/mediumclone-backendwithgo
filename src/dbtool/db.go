package dbtool

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	// dependencies for above package
	_ "github.com/jinzhu/now"
	_ "github.com/lib/pq"

	"github.com/json9512/mediumclone-backendwithgo/src/logger"
	"github.com/json9512/mediumclone-backendwithgo/src/posts"
	"github.com/json9512/mediumclone-backendwithgo/src/users"
)

// Config holds configuration for DB connection
type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUsername string
	DBPassword string
}

func getEnv(n string, dVal string) string {
	if os.Getenv(n) != "" {
		return os.Getenv(n)
	}
	return dVal
}

func createConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "mediumclone"),
		DBUsername: getEnv("DB_USERNAME", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
	}
}

// Init returns the db when connected gracefully
func Init() *gorm.DB {
	log := logger.InitLogger()
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

	return db
}

// Migrate creates necessary tables in the db
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&posts.Post{})
	db.AutoMigrate(&users.User{})
}
