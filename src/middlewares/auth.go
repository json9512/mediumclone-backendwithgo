package middlewares

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/api"
	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// VerifyUser validates the access_token in the request cookie
func VerifyUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")

		if err != nil {
			api.HandleError(c, http.StatusUnauthorized, "Token not found.")
			c.Abort()
			return
		}

		// JWT verification here
		if err := ValidateToken(c, token, pool); err != nil {
			api.HandleError(c, http.StatusUnauthorized, "Token invalid.")
			c.Abort()
			return
		}

		verifiedToken, _ := VerifyToken(token)
		username := extractUsername(verifiedToken)
		c.Set("username", username)
	}
}

// VerifyToken verifies the JWT token with the expected signing method
func VerifyToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return token, nil
}

// ValidateToken checks the validity of the provided JWT token
func ValidateToken(c context.Context, t string, pool *sql.DB) error {
	token, err := VerifyToken(t)

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return err
	}

	email, ok := claims["user_email"]

	if !ok {
		return fmt.Errorf("User email not valid.")
	}

	tokenExp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("Expiry date not valid.")
	}

	user, err := db.GetUserByEmail(c, pool, email.(string))
	if err != nil {
		return fmt.Errorf("User does not exist in DB.")
	}

	if user.TokenExpiresIn.Int64 != int64(tokenExp) {
		// This should refresh the token for the user
		return fmt.Errorf("Token expired.")
	}

	return nil
}

func extractUsername(t *jwt.Token) string {
	claims, _ := t.Claims.(jwt.MapClaims)
	email, _ := claims["user_email"]
	return strings.Split(email.(string), "@")[0]
}
