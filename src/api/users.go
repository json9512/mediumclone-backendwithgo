package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// RetrieveUser gets user by its ID from db
func RetrieveUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id := convertToInt(idStr)
		if id < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if user, err := db.GetUserByID(c, pool, id); err != nil {
			HandleError(c, http.StatusBadRequest, "User not found.")
		} else {
			c.JSON(http.StatusOK, serializeUser(user))
		}
	}
}

// RegisterUser creates a new user in db
func RegisterUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCred userInsertForm
		err := extractData(c, &userCred)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid data type.")
			return
		}

		if err := validateStruct(&userCred); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid credential.")
			return
		}
		user := bindFormToUser(&userCred)
		if user, err := db.InsertUser(c, pool, user); err != nil {
			HandleError(c, http.StatusInternalServerError, "Saving data to database failed.")
		} else {
			c.JSON(http.StatusOK, serializeUser(user))
		}
	}
}

// UpdateUser updates the user with provided info
func UpdateUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody userUpdateForm
		err := extractData(c, &reqBody)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid data type.")
			return
		}

		if err := validateStruct(reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid data.")
			return
		}

		user, err := bindUpdateFormToUser(&reqBody)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid data.")
			return
		}

		userID := int64(reqBody.ID)
		if user, err := db.UpdateUser(c, pool, userID, user); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid request.")
		} else {
			c.JSON(http.StatusOK, serializeUser(user))
		}
	}
}

// DeleteUser deletes the user in db with its ID
func DeleteUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id := convertToInt(idStr)
		if id < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if _, err := db.DeleteUserByID(c, pool, id); err != nil {
			HandleError(c, http.StatusBadRequest, "User not found.")
		} else {
			c.Status(http.StatusOK)
		}
	}
}
