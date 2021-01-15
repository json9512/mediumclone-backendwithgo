package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig ...
// Loads the configuration for the app
func LoadConfig(key string) (string, error) {
	err := godotenv.Load()

	if err != nil {
		// Check if os.Getenv(key) works
		temp := os.Getenv(key)

		if temp != "" {
			return temp, nil
		}
		return "", err
	}
	return os.Getenv(key), nil
}
