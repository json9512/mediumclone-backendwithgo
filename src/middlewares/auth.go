package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/api"
)

// VerifyUser validates the access_token in the request cookie
func VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("access_token")
		if err != nil {
			api.HandleError(c, http.StatusUnauthorized, "Unauthorized request.")
		}

		// JWT verification here
	}
}
