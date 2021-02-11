package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Logout invalidates the tokens for the user
func Logout(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInfo map[string]interface{}

		if err := c.BindJSON(&userInfo); err != nil {
			msg := "Logout failed. Invalid data type."
			HandleError(c, http.StatusBadRequest, msg)
			return
		}

		var user dbtool.User
		err := p.Query(&user, map[string]interface{}{"email": userInfo["email"]})
		if err != nil {
			msg := "Logout failed. User does not exist."
			HandleError(c, http.StatusBadRequest, msg)
			return
		}
		user.TokenExpiryAt = 0
		if err = p.Update(&user); err != nil {
			msg := "Updating user information in DB failed."
			HandleError(c, http.StatusInternalServerError, msg)
			return
		}

		c.SetCookie("access_token", "", 0, "/", "", false, true)
		c.Status(200)

	}
}
