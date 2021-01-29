package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Login validates the user and distributes the tokens
// Note: refactoring needed
func Login(p *dbtool.Pool) gin.HandlerFunc {
	handler := func(c *gin.Context) {
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

		// validate the created credential
		v := validator.New()
		if valErr := v.Struct(userCred); valErr != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. Invalid data type.",
				},
			)
			return
		}

		// Save access token and resfresh token [Testing for now]
		var user dbtool.User
		dbErr := p.Query(&user, map[string]interface{}{"email": userCred.Email})
		if dbErr != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. User does not exist.",
				},
			)
			return
		}

		// Check password
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
		user.CreatedAt = createdAt
		if err := p.Update(&user); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "Update failed.",
				},
			)
			return
		}

		c.SetCookie("access_token", "testing-access-token", 10, "/", "", false, true)
		c.Status(200)
	}
	return gin.HandlerFunc(handler)
}
