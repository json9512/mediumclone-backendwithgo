package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/api"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// VerifyUser validates the access_token in the request cookie
func VerifyUser(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")

		if err != nil {
			api.HandleError(c, http.StatusUnauthorized, "Unauthorized request. Token not found.")
			c.Abort()
			return
		}

		// JWT verification here
		if err := ValidateToken(token, p); err != nil {
			api.HandleError(c, http.StatusUnauthorized, "Unauthorized request. Token invalid.")
			c.Abort()
			return
		}

	}
}

func VerifyToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ValidateToken(t string, p *dbtool.Pool) error {
	token, err := VerifyToken(t)

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return err
	}

	_, ok = claims["user_email"]

	if !ok {
		return fmt.Errorf("User email not valid")
	}

	tokenExp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("Expiry date not valid")
	}

	user := &dbtool.User{}
	if err := p.Query(&user, nil); err != nil {
		return fmt.Errorf("User does not exist in DB")
	}

	if user.TokenExpiryAt != int64(tokenExp) {
		return fmt.Errorf("Token expired")
	}

	return nil
}
