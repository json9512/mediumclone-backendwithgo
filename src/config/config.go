package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// LoadConfig ...
// Loads the configuration for the app
func LoadConfig(key string) string {

	readVariablesFromFile(".env")

	// Check if os.Getenv(key) has value
	temp := os.Getenv(key)
	if temp != "" {
		return temp
	}

	log.Fatal("Failed to load variable")
	return os.Getenv(key)
}

func readVariablesFromFile(filename string) {
	envFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(envFile)

	// Add key=value as environment variables
	for {
		dataRead, err := reader.ReadBytes('\n')
		dataString := string(dataRead)
		dataSplit := strings.Split(dataString, "=")

		os.Setenv(dataSplit[0], dataSplit[1])

		// Exit when nothing to read
		if err != nil {
			break
		}
	}
}
