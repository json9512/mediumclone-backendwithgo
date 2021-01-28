package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Login validates the user and distributes the tokens
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

		// Save access token and resfresh token [Testing for now]
		var user dbtool.User
		dbQuery := p.Where("email = ?", userCred.Email).Find(&user)
		if dbQuery.Error != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Authentication failed. User does not exist.",
				},
			)
			return
		}

		// Update the TokenCreatedAt time
		createdAt := time.Now()
		p.Model(&user).Updates(
			map[string]interface{}{
				"token_created_at": createdAt,
			})

		c.SetCookie("access_token", "testing-access-token", 10, "/", "", false, true)
		c.Status(200)
	}
	return gin.HandlerFunc(handler)
}
