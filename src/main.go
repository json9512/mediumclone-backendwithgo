//Package main implements the Go RESTful API for mediumclone project
package main

import (
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
	"github.com/json9512/mediumclone-backendwithgo/src/logger"
	"github.com/json9512/mediumclone-backendwithgo/src/posts"
	"github.com/json9512/mediumclone-backendwithgo/src/users"
)

// SetupRouter returns a *gin.Engine
func SetupRouter(mode string) *gin.Engine {
	var router *gin.Engine
	log := logger.InitLogger()
	config.ReadVariablesFromFile(".env")

	if mode != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()

	if mode != "test" {
		router.Use(logger.Middleware(log))
	}

	router.Use(gin.Recovery())

	// Add routes
	posts.AddRoutes(router)
	users.AddRoutes(router)
	return router
}

func main() {

	r := SetupRouter("debug")

	// db connection
	db := dbtool.Init()
	dbtool.Migrate(db)
	defer db.Close()

	r.Run() // Port 8080
}
