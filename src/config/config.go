package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig ...
// Loads the configuration for the app
func LoadConfig(key string) string {
	err := godotenv.Load()

	if err != nil {
		// Check if os.Getenv(key) has value
		temp := os.Getenv(key)

		if temp != "" {
			return temp
		}

		log.Fatal("godotenv failed to load variable", err)
	}
	return os.Getenv(key)
}
