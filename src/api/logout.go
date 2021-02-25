package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Logout invalidates the tokens for the user
func Logout(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInfo map[string]interface{}

		if err := c.BindJSON(&userInfo); err != nil {
			HandleError(c, http.StatusBadRequest, "Logout failed. Invalid data type.")
			return
		}

		email := userInfo["email"]
		user, err := db.GetUserByEmail(email.(string))
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Logout failed. User does not exist.")
			return
		}

		query, err := createUpdateQuery(user.ID, user.Email, user.Password, 0)
		if _, err = db.UpdateUser(query); err != nil {
			HandleError(c, http.StatusInternalServerError, "Updating user information in DB failed.")
			return
		}

		c.SetCookie("access_token", "", 0, "/", "", false, true)
		c.Status(200)

	}
}
