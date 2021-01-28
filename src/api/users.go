package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// RetrieveUser gets user by its ID from db
func RetrieveUser(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Invalid ID.",
				},
			)
			return
		}

		var user dbtool.User
		dbErr := p.Query(&user, map[string]interface{}{"id": idInt})
		if dbErr != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "User not found",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			serializeUser(user),
		)
		return
	}
}

// RegisterUser creates a new user in db
func RegisterUser(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// NOTE to future me:
		// - need to implement input verification
		// - refactoring needed
		var reqBody UserReqData
		err := handleReqBody(c, &reqBody, "User registration failed. Invalid data type.")
		if err != nil {
			return
		}

		// Convert reqBody to User type with empty access token and refresh token
		user := createUserObj(reqBody)

		// Save to db
		dbErr := p.Insert(&user)
		if dbErr != nil {
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "User registration failed. Saving data to database failed.",
				},
			)
			return
		}
		// Serialize data
		c.JSON(
			http.StatusOK,
			serializeUser(user))
	}
}

// UpdateUser updates the user with provided info
func UpdateUser(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody UserReqData
		err := handleReqBody(c, &reqBody, "User update failed. Invalid data type.")
		if err != nil {
			return
		}
		// NOTE: need to retreive access token and refresh token from header
		user := createUserObj(reqBody)

		dbErr := p.Update(&user)
		if dbErr != nil {
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "User update failed. Saving data to database failed.",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			serializeUser(user),
		)
	}
}

// DeleteUser deletes the user in db with its ID
func DeleteUser(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Invalid ID",
				},
			)
			return
		}

		var user dbtool.User
		dbErr := p.Delete(&user, map[string]interface{}{"id": idInt})
		if dbErr != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Deleting user data from database failed. User not found",
				},
			)
			return
		}

		c.Status(http.StatusOK)
	}
}

func handleReqBody(c *gin.Context, reqBody *UserReqData, errorMsg string) error {
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(
			http.StatusBadRequest,
			&errorResponse{
				Msg: errorMsg,
			},
		)
		return err
	}
	return nil
}
