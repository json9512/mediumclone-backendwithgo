package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Login validates the user and distributes the tokens
func Login(db *dbtool.DB, envVars *config.EnvVars) gin.HandlerFunc {
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

		user, err := db.GetUserByEmail(userCred.Email)
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

		expiresIn := time.Now().Add(time.Hour * 24).Unix()
		user.TokenExpiresIn = expiresIn
		if err := db.Update(&user); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "Authentication failed. Update failed.",
				},
			)
			return
		}

		at, err := createAccessToken(user.Email, envVars.JWTSecret, expiresIn)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, "Login failed. Unable to create token.")
		}

		c.SetCookie("access_token", at, 10, "/", "", false, true)
		c.Status(200)
	}
}

func createAccessToken(userEmail, secret string, expiryDate int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = userEmail
	claims["exp"] = expiryDate
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := accessToken.SignedString([]byte(secret))
	return token, err
}
