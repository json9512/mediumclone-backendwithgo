package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Logout invalidates the tokens for the user
func Logout(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var user dbtool.User
		err := p.Query(&user, map[string]interface{}{"email": userInfo["email"]})
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Logout failed. User does not exist.",
				},
			)
			return
		}
		user.TokenCreatedAt = &time.Time{}
		err = p.Update(&user)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "Updating user information in DB failed.",
				},
			)
			return
		}

		c.SetCookie("access_token", "", 0, "/", "", false, true)
		c.Status(200)

	}
}
