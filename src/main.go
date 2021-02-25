//Package main implements the Go RESTful API for mediumclone project
package main

import (
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
	"github.com/json9512/mediumclone-backendwithgo/src/middlewares"
	"github.com/json9512/mediumclone-backendwithgo/src/routes"
)

// SetupRouter returns the API server
func SetupRouter(mode string, db *dbtool.DB) *gin.Engine {
	var router *gin.Engine
	log := config.InitLogger()
	envVars := config.LoadEnvVars()

	if mode != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()

	if mode != "test" {
		router.Use(middlewares.CustomLogger(log))
	}

	router.Use(gin.Recovery())
	routes.AddRoutes(router, db, &envVars)
	return router
}

func main() {
	config.ReadVariablesFromFile(".env")
	// db connection
	db := dbtool.Init()
	dbtool.Migrate(db)
	defer db.Close()

	r := SetupRouter("debug", db)

	r.Run() // Port 8080
}
