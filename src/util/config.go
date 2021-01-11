package util

import (
	"github.com/spf13/viper"
)

// Config ...
// holds configuration for the app
type Config struct {
	Database DatabaseConfig
}

// DatabaseConfig ...
type DatabaseConfig struct {
	DBUsername string
	DBPassword string
	DBPort     string
	DBHost     string
	DBName     string
}

// LoadConfig ...
// Loads the configuration for the app
func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("viper")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}