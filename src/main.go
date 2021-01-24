//Package main implements the Go RESTful API for mediumclone project
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
	"github.com/json9512/mediumclone-backendwithgo/src/logger"
	"github.com/json9512/mediumclone-backendwithgo/src/login"
	"github.com/json9512/mediumclone-backendwithgo/src/posts"
	"github.com/json9512/mediumclone-backendwithgo/src/users"
)

// SetupRouter returns a *gin.Engine
func SetupRouter(mode string, db *gorm.DB) *gin.Engine {
	var router *gin.Engine
	log := logger.InitLogger()

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
	users.AddRoutes(router, db)
	login.AddRoutes(router, db)
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
