package config

import (
	"bufio"
	"os"
	"strings"

	"github.com/json9512/mediumclone-backendwithgo/src/logger"
)

// ReadVariablesFromFile ...
// reads environment variables from given filename
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
