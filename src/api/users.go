package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// RetrieveUser godoc
// @Summary Get user
// @Tags users
// @Description Get user by its ID
// @ID get-user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} api.SwaggerUser
// @Failure 400 {object} api.APIError "Bad Request"
// @Router /users/:id [get]
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

// RegisterUser godoc
// @Summary Create new user
// @Tags users
// @Description Create a new user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param userInfo body api.UserInsertForm true "Add user"
// @Success 200 {object} api.SwaggerUser
// @Failure 400 {object} api.APIError "Bad Request"
// @Failure 500 {object} api.APIError "Internal Server Error"
// @Failure 401 {object} api.APIError "Unauthorized"
// @Router /users [post]
func RegisterUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCred UserInsertForm
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

// UpdateUser godoc
// @Summary Update user
// @Tags users
// @Description Update user with provided information
// @ID update-user
// @Accept  json
// @Produce  json
// @Param userInfo body api.UserUpdateForm true "Update user"
// @Success 200 {object} api.SwaggerUser
// @Failure 400 {object} api.APIError "Bad Request"
// @Failure 401 {object} api.APIError "Unauthorized"
// @Router /users [put]
func UpdateUser(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody UserUpdateForm
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

// DeleteUser godoc
// @Summary Delete user
// @Tags users
// @Description Delete user by its ID
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param id path int true "Delete user"
// @Success 200
// @Failure 400 {object} api.APIError "Bad Request"
// @Failure 401 {object} api.APIError "Unauthorized"
// @Router /users [delete]
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
