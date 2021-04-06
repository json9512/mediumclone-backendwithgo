package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// Logout invalidates the tokens for the user
func Logout(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInfo map[string]interface{}

		if err := c.BindJSON(&userInfo); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid data type.")
			return
		}

		email := userInfo["email"]
		user, err := db.GetUserByEmail(c, pool, email.(string))
		if err != nil {
			HandleError(c, http.StatusBadRequest, "User does not exist.")
			return
		}

		user, err = db.UpdateTokenExpiresIn(c, pool, user, 0)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, "Updating user information in DB failed.")
			return
		}

		c.SetCookie("access_token", "", 0, "/", "", false, true)
		c.Status(200)

	}
}
