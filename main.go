package main

import "github.com/gin-gonic/gin"

// SetupRouter ...
// returns a *gin.Engine
func SetupRouter(mode string) *gin.Engine {
	var router *gin.Engine

	// Set gin mode and create router
	if mode != "debug" {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		router.Use(gin.Recovery())
	} else {
		// Append logger and recovery middleware if debug mode
		router = gin.New()
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
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
