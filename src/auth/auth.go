package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/json9512/mediumclone-backendwithgo/src/users"
)

type loginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type errorResponse struct {
	Msg string `json:"message"`
}

// AddRoutes adds HTTP Methods for the /posts endpoint
func AddRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/login", func(c *gin.Context) {
		var userCredentials loginCredentials
		if err := c.BindJSON(&userCredentials); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. Invalid data type.",
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
				&errorResponse{
					Msg: "Authentication failed. User does not exist.",
				},
			)
			return
		}

		db.Model(&user).Updates(
			map[string]interface{}{
				"access_token":  "testing-access-token",
				"refresh_token": "testing-refresh-token",
			})

		c.SetCookie("access_token", "testing-access-token", 10, "/", "", false, true)
		c.SetCookie("refresh_token", "testing-refresh-token", 10, "/", "", false, true)

		c.Status(200)
	})

	router.POST("/logout", func(c *gin.Context) {
		var userDetail map[string]interface{}
		if err := c.BindJSON(&userDetail); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Logout failed. Invalid data type.",
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
				&errorResponse{
					Msg: "Authentication failed. User does not exist.",
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

		c.SetCookie("access_token", "", 0, "/", "", false, true)
		c.SetCookie("refresh_token", "", 0, "/", "", false, true)

		c.Status(200)

	})

}
