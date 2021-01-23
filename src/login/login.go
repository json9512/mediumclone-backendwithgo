package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/json9512/mediumclone-backendwithgo/src/config"
)

type loginCredentials struct {
	email    string
	password string
}

// AddRoutes adds HTTP Methods for the /posts endpoint
func AddRoutes(router *gin.Engine) {
	router.POST("/login", func(c *gin.Context) {
		var userCredentials loginCredentials
		err := c.BindJSON(&userCredentials)
		if err != nil {
			errResponse := &config.ResData{
				"message": "Login failed with username and password",
			}
			c.JSON(
				http.StatusBadRequest,
				errResponse,
			)
		}

		c.JSON(200, &config.ResData{
			"access-token":  "testing-access-token",
			"refresh-token": "testing-refresh-token",
		})
	})

}
