// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

package main

//go:generate sqlboiler --wipe psql

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/json9512/mediumclone-backendwithgo/src/docs"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

func main() {
	logger := config.InitLogger()
	config.ReadVariablesFromFile(".env")

	dbContainer := DBProvider.Init(logger)
	err := dbContainer.Migrate("up")
	if err != nil {
		logger.Error(err)
	}

	r := SetupRouter("debug", logger, dbContainer.DB)
	r.Run() // Port 8080
}
