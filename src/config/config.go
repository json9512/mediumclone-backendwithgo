package config

import (
	"bufio"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// EnvVars holds environment variables necessary for the server
type EnvVars struct {
	JWTSecret string
}

// InitLogger returns a formatted logger
func InitLogger() *logrus.Logger {
	log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}

// ReadVariablesFromFile reads environment variables from given filename
func ReadVariablesFromFile(filename string) {
	log := InitLogger()

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

// LoadEnvVars load environment variables necessary for the server
func LoadEnvVars() *EnvVars {
	return &EnvVars{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

}
