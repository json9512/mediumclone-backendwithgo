//Package main implements the Go RESTful API for mediumclone project
package main

import (
	"fmt"
	dbManager "json9512/mediumclone-go/db"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
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

		log.SetFormatter(&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
		})

		router.Use(ginlogrus.Logger(log))
		router.Use(gin.Recovery())

		// Check for db
		db, msg, err := dbManager.ConnectDB()
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
	fmt.Println("Test CODEOWNERS")
}
