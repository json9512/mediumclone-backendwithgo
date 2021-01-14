package util

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig ...
// Loads the configuration for the app
func LoadConfig(key string) (string, error) {
	err := godotenv.Load()

	if err != nil {
		return "", err
	}

	return os.Getenv(key), nil
}
