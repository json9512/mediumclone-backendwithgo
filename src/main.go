//Package main implements the Go RESTful API for mediumclone project
package main

import (
	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// SetupRouter ...
// returns a *gin.Engine
func SetupRouter(mode string) *gin.Engine {
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

		// Check for db
		db, msg, err := db.ConnectDB()
		if db == nil || err != nil {
			log.WithField("Error", err).Fatal(msg)
		} else {
			log.Info(msg)
		}
		defer db.Close()
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}

func main() {
	r := SetupRouter("debug")
	r.Run() // Port 8080
}
