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
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if user, err := db.GetUserByID(idInt); err != nil {
			HandleError(c, http.StatusBadRequest, "User not found.")
		} else {
			c.JSON(http.StatusOK, serializeUser(user))
		}
	}
}

// RegisterUser creates a new user in db
func RegisterUser(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCred credential
		err := extractData(c, &userCred)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "User registration failed. Invalid data type.")
			return
		}

		if err := validateStruct(&userCred); err != nil {
			HandleError(c, http.StatusBadRequest, "User registration failed. Invalid credential.")
			return
		}

		if user, err := db.CreateUser(userCred.Email, userCred.Password); err != nil {
			HandleError(c, http.StatusInternalServerError, "User update failed. Saving data to database failed.")
		} else {
			c.JSON(http.StatusOK, serializeUser(user))
		}
	}
}

// UpdateUser updates the user with provided info
func UpdateUser(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody userUpdateForm
		err := extractData(c, &reqBody)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "User update failed. Invalid data type.")
			return
		}

		if err := validateStruct(reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "User update failed. Invalid data.")
			return
		}

		query, err := createUserUpdateQuery(reqBody.ID, reqBody.Email, reqBody.Password, nil)
		if err != nil {
			HandleError(c, http.StatusBadRequest, err.Error())
			return
		}

		if !db.CheckIfUserExists(query.ID) {
			HandleError(c, http.StatusBadRequest, "User update failed. Invalid ID.")
			return
		}

		if updatedUser, err := db.UpdateUser(&query); err != nil {
			HandleError(c, http.StatusInternalServerError, "User update failed. Saving data to database failed.")
		} else {
			c.JSON(http.StatusOK, serializeUser(updatedUser))
		}
	}
}

// DeleteUser deletes the user in db with its ID
func DeleteUser(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)

		if err != nil || idInt < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		if _, err := db.DeleteUserByID(idInt); err != nil {
			HandleError(c, http.StatusBadRequest, "Deleting user data from database failed. User not found")
		} else {
			c.Status(http.StatusOK)
		}

	}
}

func extractData(c *gin.Context, reqBody interface{}) error {
	if err := c.BindJSON(&reqBody); err != nil {
		return err
	}
	return nil
}
