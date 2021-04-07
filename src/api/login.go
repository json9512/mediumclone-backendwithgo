package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

func Login(pool *sql.DB, env *config.EnvVars) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCred userInsertForm
		if err := c.BindJSON(&userCred); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&APIError{
					Msg: "Invalid data type.",
				},
			)
			return
		}

		if err := validateStruct(&userCred); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&APIError{
					Msg: "Invalid data type.",
				},
			)
			return
		}

		user, err := db.GetUserByEmail(c, pool, userCred.Email)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&APIError{
					Msg: "User does not exist.",
				},
			)
			return
		}

		if user.PWD.String != userCred.Password {
			c.JSON(
				http.StatusBadRequest,
				&APIError{
					Msg: "Wrong password.",
				},
			)
			return
		}

		expiresIn := time.Now().Add(time.Hour * 24).Unix()
		user, err = db.UpdateTokenExpiresIn(c, pool, user, expiresIn)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				&APIError{
					Msg: "Authentication failed. Update failed.",
				},
			)
			return
		}

		at, err := CreateAccessToken(user.Email.String, env.JWTSecret, expiresIn)
		if err != nil {
			HandleError(c, http.StatusInternalServerError, "Login failed. Unable to create token.")
		}

		c.SetCookie("access_token", at, 10, "/", "", false, true)
		c.Status(200)
	}
}

func CreateAccessToken(userEmail, secret string, expiryDate int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = userEmail
	claims["exp"] = expiryDate
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := accessToken.SignedString([]byte(secret))
	return token, err
}
