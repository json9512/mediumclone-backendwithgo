package api

import (
	"fmt"
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
					Msg: "User not found.",
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
		var userCred credential
		err := handleReqBody(c, &userCred, "User registration failed. Invalid data type.")
		if err != nil {
			return
		}

		if valErr := validateStruct(&userCred); valErr != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "User registration failed. Invalid credential.",
				},
			)
			return
		}

		user := dbtool.User{
			Email:    userCred.Email,
			Password: userCred.Password,
		}

		dbErr := p.Insert(&user)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "User registration failed. Saving data to database failed.",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			serializeUser(user))
	}
}

// UpdateUser updates the user with provided info
func UpdateUser(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody userUpdateForm
		err := handleReqBody(c, &reqBody, "User update failed. Invalid data type.")
		if err != nil {
			return
		}

		if valErr := validateStruct(reqBody); valErr != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "User update failed. Invalid data.",
				},
			)
			return
		}

		user, err := createUserUpdate(reqBody)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: err.Error(),
				},
			)
			return
		}

		qErr := p.Query(&dbtool.User{}, map[string]interface{}{"id": user.ID})
		if qErr != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "User update failed. Invalid ID.",
				},
			)
			return
		}

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

func handleReqBody(c *gin.Context, reqBody interface{}, errorMsg string) error {
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
