package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// RetrieveUser gets user by its ID from db
func RetrieveUser(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			msg := "Invalid ID."
			HandleError(c, http.StatusBadRequest, msg)
			return
		}

		user, err := db.GetUserByID(idInt)
		if err != nil {
			msg := "User not found."
			HandleError(c, http.StatusBadRequest, msg)
			return
		}

		c.JSON(
			http.StatusOK,
			serializeUser(user),
		)
	}
}

// RegisterUser creates a new user in db
func RegisterUser(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCred credential
		err := extractData(c, &userCred, "User registration failed. Invalid data type.")
		if err != nil {
			return
		}

		if err := validateStruct(&userCred); err != nil {
			msg := "User registration failed. Invalid credential."
			HandleError(c, http.StatusBadRequest, msg)
			return
		}

		user := dbtool.User{
			Email:    userCred.Email,
			Password: userCred.Password,
		}

		if err := db.Insert(&user); err != nil {
			msg := "User update failed. Saving data to database failed."
			HandleError(c, http.StatusInternalServerError, msg)
			return
		}

		c.JSON(
			http.StatusOK,
			serializeUser(&user))
	}
}

// UpdateUser updates the user with provided info
func UpdateUser(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody userUpdateForm
		err := extractData(c, &reqBody, "User update failed. Invalid data type.")
		if err != nil {
			return
		}

		if err := validateStruct(reqBody); err != nil {
			msg := "User update failed. Invalid data."
			HandleError(c, http.StatusBadRequest, msg)
			return
		}

		user, err := createUserUpdate(reqBody)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err.Error())
			return
		}

		err = db.Query(&dbtool.User{}, map[string]interface{}{"id": user.ID})
		if err != nil {
			msg := "User update failed. Invalid ID."
			HandleError(c, http.StatusBadRequest, msg)
			return
		}

		if err = db.Update(&user); err != nil {
			msg := "User update failed. Saving data to database failed."
			HandleError(c, http.StatusInternalServerError, msg)
			return
		}

		c.JSON(http.StatusOK, serializeUser(&user))
	}
}

// DeleteUser deletes the user in db with its ID
func DeleteUser(p *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)

		if err != nil || idInt < 1 {
			msg := "Invalid ID"
			HandleError(c, http.StatusBadRequest, msg)
			return
		}

		var user dbtool.User
		err = p.Delete(&user, map[string]interface{}{"id": idInt})
		if err != nil {
			msg := "Deleting user data from database failed. User not found"
			HandleError(c, http.StatusBadRequest, msg)
			return
		}

		c.Status(http.StatusOK)
	}
}

func extractData(c *gin.Context, reqBody interface{}, errorMsg string) error {
	if err := c.BindJSON(&reqBody); err != nil {
		HandleError(c, http.StatusBadRequest, errorMsg)
		return err
	}
	return nil
}
