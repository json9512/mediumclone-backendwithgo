//Package main implements the Go RESTful API for mediumclone project
package main

//go:generate sqlboiler --wipe psql

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	DBProvider "github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/middlewares"
	"github.com/json9512/mediumclone-backendwithgo/src/routes"
)

// SetupRouter returns the API server
func SetupRouter(mode string, logger *logrus.Logger, db *sql.DB) *gin.Engine {
	var router *gin.Engine
	envVars := config.LoadEnvVars()

	if mode != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()
	if mode != "test" {
		router.Use(middlewares.CustomLogger(logger))
	}

	router.Use(gin.Recovery())
	routes.AddRoutes(router, db, envVars)
	return router
}

func main() {
	logger := config.InitLogger()
	config.ReadVariablesFromFile(".env")

	dbContainer := DBProvider.Init(logger)
	err := dbContainer.Migrate(logger, "up")
	if err != nil {
		logger.Error(err)
	}

	r := SetupRouter("debug", logger, dbContainer.DB)
	r.Run() // Port 8080
}
