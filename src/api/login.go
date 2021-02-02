package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Login validates the user and distributes the tokens
// Note: refactoring needed
func Login(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCred credential
		if err := c.BindJSON(&userCred); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. Invalid data type.",
				},
			)
			return
		}

		if err := validateStruct(&userCred); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. Invalid data type.",
				},
			)
			return
		}

		var user dbtool.User
		err := p.Query(&user, map[string]interface{}{"email": userCred.Email})
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. User does not exist.",
				},
			)
			return
		}

		if user.Password != userCred.Password {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. Wrong password.",
				},
			)
			return
		}

		// Update the TokenCreatedAt time
		createdAt := time.Now()
		user.TokenCreatedAt = &createdAt
		if err := p.Update(&user); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "Authentication failed. Update failed.",
				},
			)
			return
		}

		c.SetCookie("access_token", "testing-access-token", 10, "/", "", false, true)
		c.Status(200)
	}
}
