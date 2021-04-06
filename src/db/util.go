package db

import (
	"database/sql/driver"
	"encoding/json"
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

// JSONB type def for database
type JSONB map[string]interface{}

// Value for gorm to read the JSONB data
func (j JSONB) Value() (driver.Value, error) {
	valString, err := json.Marshal(j)
	return string(valString), err
}

// Scan for gorm to scan the JSONB data
func (j *JSONB) Scan(v interface{}) error {
	err := json.Unmarshal(v.([]byte), &j)
	if err != nil {
		return err
	}
	return nil
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
