package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RetrieveUser gets user by its ID from db
func RetrieveUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if user, err := pool.GetUserByID(idInt); err != nil {
			HandleError(c, http.StatusBadRequest, "User not found.")
		} else {
			c.JSON(http.StatusOK, serializeUser(user))
		}
	}
}

// RegisterUser creates a new user in db
func RegisterUser(pool *sql.DB) gin.HandlerFunc {
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

		if user, err := pool.CreateUser(userCred.Email, userCred.Password); err != nil {
			HandleError(c, http.StatusInternalServerError, "User update failed. Saving data to database failed.")
		} else {
			c.JSON(http.StatusOK, serializeUser(user))
		}
	}
}

// UpdateUser updates the user with provided info
func UpdateUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserForm
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

		if !pool.CheckIfUserExists(query.ID) {
			HandleError(c, http.StatusBadRequest, "User update failed. Invalid ID.")
			return
		}

		if _, err := pool.UpdateUser(&query); err != nil {
			HandleError(c, http.StatusInternalServerError, "User update failed. Saving data to database failed.")
		} else {
			c.Status(http.StatusOK)
		}
	}
}

// DeleteUser deletes the user in db with its ID
func DeleteUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)

		if err != nil || idInt < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		if _, err := pool.DeleteUserByID(idInt); err != nil {
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
