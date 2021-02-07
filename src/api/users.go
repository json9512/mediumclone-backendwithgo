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
			msg := "Invalid ID."
			handleError(&customError{c, http.StatusBadRequest, msg})
			return
		}

		var user dbtool.User
		err = p.Query(&user, map[string]interface{}{"id": idInt})
		if err != nil {
			msg := "User not found."
			handleError(&customError{c, http.StatusBadRequest, msg})
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

		if err := validateStruct(&userCred); err != nil {
			msg := "User registration failed. Invalid credential."
			handleError(&customError{c, http.StatusBadRequest, msg})
			return
		}

		user := dbtool.User{
			Email:    userCred.Email,
			Password: userCred.Password,
		}

		if err := p.Insert(&user); err != nil {
			msg := "User update failed. Saving data to database failed."
			handleError(&customError{c, http.StatusInternalServerError, msg})
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

		if err := validateStruct(reqBody); err != nil {
			msg := "User update failed. Invalid data."
			handleError(&customError{c, http.StatusBadRequest, msg})
			return
		}

		user, err := createUserUpdate(reqBody)
		if err != nil {
			handleError(&customError{c, http.StatusBadRequest, err.Error()})
			return
		}

		err = p.Query(&dbtool.User{}, map[string]interface{}{"id": user.ID})
		if err != nil {
			msg := "User update failed. Invalid ID."
			handleError(&customError{c, http.StatusBadRequest, msg})
			return
		}

		if err = p.Update(&user); err != nil {
			msg := "User update failed. Saving data to database failed."
			handleError(&customError{c, http.StatusInternalServerError, msg})
			return
		}

		c.JSON(http.StatusOK, serializeUser(user))
	}
}

// DeleteUser deletes the user in db with its ID
func DeleteUser(p *dbtool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)

		if err != nil || idInt < 1 {
			msg := "Invalid ID"
			handleError(&customError{c, http.StatusBadRequest, msg})
			return
		}

		var user dbtool.User
		err = p.Delete(&user, map[string]interface{}{"id": idInt})
		if err != nil {
			msg := "Deleting user data from database failed. User not found"
			handleError(&customError{c, http.StatusBadRequest, msg})
			return
		}

		c.Status(http.StatusOK)
	}
}

func handleReqBody(c *gin.Context, reqBody interface{}, errorMsg string) error {
	if err := c.BindJSON(&reqBody); err != nil {
		handleError(&customError{c, http.StatusBadRequest, errorMsg})
		return err
	}
	return nil
}
