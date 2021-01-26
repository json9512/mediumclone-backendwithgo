package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Logout invalidates the tokens for the user
func Logout(p *dbtool.Pool) gin.HandlerFunc {
	handler := func(c *gin.Context) {
		var userInfo map[string]interface{}
		if err := c.BindJSON(&userInfo); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Logout failed. Invalid data type.",
				},
			)
			return
		}

		// Save access token and resfresh token [Testing for now]
		var user dbtool.User
		dbQuery := p.Where("email = ?", userInfo["email"]).Find(&user)
		if dbQuery.Error != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. User does not exist.",
				},
			)
			return
		}

		// Update the database
		p.Model(&user).Updates(
			map[string]interface{}{
				"access_token":  "",
				"refresh_token": "",
			})

		c.SetCookie("access_token", "", 0, "/", "", false, true)
		c.SetCookie("refresh_token", "", 0, "/", "", false, true)

		c.Status(200)

	}

	return gin.HandlerFunc(handler)
}
