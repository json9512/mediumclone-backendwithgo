package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/api"
)

// VerifyUser validates the access_token in the request cookie
func VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken, err := c.Cookie("access_token")
		fmt.Println("user token from middleware", userToken)
		if err != nil {
			api.HandleError(&api.CustomError{
				G:    c,
				Code: http.StatusUnauthorized,
				Msg:  "Unauthorized request."})
		}

		// JWT verification here
	}
}
