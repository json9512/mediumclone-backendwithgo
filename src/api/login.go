package api

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Login validates the user and distributes the tokens
// Note: refactoring needed
func Login(db *dbtool.DB) gin.HandlerFunc {
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
		err := db.Query(&user, map[string]interface{}{"email": userCred.Email})
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
		expiryDate := time.Now().Add(time.Hour * 24).Unix()
		user.TokenExpiryDate = expiryDate
		if err := db.Update(&user); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "Authentication failed. Update failed.",
				},
			)
			return
		}

		at, err := createAccessToken(user.Email, expiryDate)

		if err != nil {
			msg := "Login failed. Unable to create token."
			HandleError(c, http.StatusInternalServerError, msg)
		}

		c.SetCookie("access_token", at, 10, "/", "", false, true)
		c.Status(200)
	}
}

func createAccessToken(userEmail string, expiryDate int64) (string, error) {
	var err error
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = userEmail
	claims["exp"] = expiryDate
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
