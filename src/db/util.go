package db

import (
	"os"
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
