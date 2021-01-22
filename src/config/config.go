package config

import (
	"bufio"
	"database/sql/driver"
	"encoding/json"
	"os"
	"strings"

	"github.com/json9512/mediumclone-backendwithgo/src/logger"
)

// JSONB type def for database
type JSONB map[string]interface{}

// ReadVariablesFromFile reads environment variables from given filename
func ReadVariablesFromFile(filename string) {
	log := logger.InitLogger()

	// Github Actions 테스트 환경에서는 .env파일이 없다
	envFile, err := os.Open(filename)
	if err != nil {
		log.Info(err)
		return
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
