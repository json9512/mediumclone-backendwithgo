//Package main implements the Go RESTful API for mediumclone project
package main

import (
	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/posts"
)

// SetupRouter ...
// returns a *gin.Engine
func SetupRouter(mode string) (*gin.Engine, *logrus.Logger) {
	var router *gin.Engine
	log := logrus.New()

	// Set gin mode and create router
	if mode != "debug" {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		router.Use(gin.Recovery())
	} else {
		// Append logger and recovery middleware if debug mode
		router = gin.New()

		log.SetFormatter(&formatter.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
		})

		router.Use(ginlogrus.Logger(log))
		router.Use(gin.Recovery())
	}

	// For test only
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Add routes
	posts.AddRoutes(router)

	return router, log
}

func main() {

	r, logger := SetupRouter("debug")

	// db connection
	db := db.Init()
	defer db.Close()

	logger.Info("DB connection successful")

	r.Run() // Port 8080
}
