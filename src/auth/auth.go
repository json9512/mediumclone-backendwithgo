package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/users"
)

type loginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AddRoutes adds HTTP Methods for the /posts endpoint
func AddRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/login", func(c *gin.Context) {
		var userCredentials loginCredentials
		err := c.BindJSON(&userCredentials)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&config.ResData{
					"message": "Authentication failed. Invalid data type.",
				},
			)
			return
		}

		// Save access token and resfresh token [Testing for now]
		var user users.User
		dbResult := db.Where("email = ?", userCredentials.Email).Find(&user)
		if dbResult.Error != nil {
			c.JSON(
				http.StatusBadRequest,
				&config.ResData{
					"message": "Authentication failed. User does not exist.",
				},
			)
			return
		}

		// Update the database
		db.Model(&user).Updates(
			map[string]interface{}{
				"access_token":  "testing-access-token",
				"refresh_token": "testing-refresh-token",
			})

		c.JSON(200, &config.ResData{
			"access-token":  "testing-access-token",
			"refresh-token": "testing-refresh-token",
		})
	})

	router.POST("/logout", func(c *gin.Context) {
		var userDetail map[string]interface{}
		err := c.BindJSON(&userDetail)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&config.ResData{
					"message": "Logout failed. Invalid data type.",
				},
			)
			return
		}

		// Save access token and resfresh token [Testing for now]
		var user users.User
		dbResult := db.Where("email = ?", userDetail["email"]).Find(&user)
		if dbResult.Error != nil {
			c.JSON(
				http.StatusBadRequest,
				&config.ResData{
					"message": "Authentication failed. User does not exist.",
				},
			)
			return
		}

		// Update the database
		db.Model(&user).Updates(
			map[string]interface{}{
				"access_token":  "",
				"refresh_token": "",
			})

		c.JSON(200, &config.ResData{
			"access-token":  "",
			"refresh-token": "",
		})

	})

}
