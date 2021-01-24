package users

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type userReqData struct {
	UserID   uint   `json:"user-id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AddRoutes adds HTTP Methods for the /users endpoint
func AddRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/users/:id", retrieveUser(db))

	router.POST("/users", registerUser(db))

	router.PUT("/users", updateUser(db))

	router.DELETE("/users/:id", deleteUser(db))
}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}

func retrieveUser(db *gorm.DB) gin.HandlerFunc {
	handler := func(c *gin.Context) {
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

		var user User
		result := db.Find(&user, idInt)
		if result.Error != nil {
			fmt.Println(result.Error, id)
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
			serialize(&user),
		)
		return
	}

	return gin.HandlerFunc(handler)
}

func registerUser(db *gorm.DB) gin.HandlerFunc {
	handler := func(c *gin.Context) {
		// NOTE to future me:
		// - need to implement input verification
		// - refactoring needed
		var reqBody userReqData
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "User registration failed. Invalid data type.",
				},
			)
			return
		}
		// Convert reqBody to User type with empty access token and refresh token
		user := CreateUserData(reqBody, "", "")

		// Save to db
		result := db.Create(&user)
		if result.Error != nil {
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
			serialize(user))
	}
	return gin.HandlerFunc(handler)
}

func updateUser(db *gorm.DB) gin.HandlerFunc {
	handler := func(c *gin.Context) {
		var reqBody userReqData
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "User update failed. Invalid data type.",
				},
			)
			return
		}
		// NOTE: need to retreive access token and refresh token from header
		user := CreateUserData(reqBody, "", "")

		result := db.Model(&User{}).Updates(user)
		if result.Error != nil {
			c.JSON(
				http.StatusInternalServerError,
				&errorResponse{
					Msg: "User Update failed. Saving data to database failed.",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			serialize(user),
		)
	}
	return gin.HandlerFunc(handler)
}

func deleteUser(db *gorm.DB) gin.HandlerFunc {
	handler := func(c *gin.Context) {
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

		// NOTE: need to retreive access token and refresh token from header
		var user User
		result := db.Find(&user, idInt).Delete(user)
		if result.Error != nil {
			c.JSON(
				http.StatusBadRequest,
				&errorResponse{
					Msg: "Deleting user data from database failed. User not found",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			serialize(&user),
		)
	}
	return gin.HandlerFunc(handler)
}
