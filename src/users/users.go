package users

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type resData map[string]interface{}

type userReqData struct {
	UserID string `json:"user-id"`
	Email  string `json:"email"`
}

// AddRoutes adds HTTP Methods for the /users endpoint
func AddRoutes(router *gin.Engine) {
	router.GET("/users", func(c *gin.Context) {

		// Check queries
		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
			c.JSON(200, resData{
				"result": queries,
			})
		} else {
			c.JSON(200, resData{
				"result": []string{"test", "sample", "users"},
			})
		}
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, resData{
			"result": id,
		})
	})

	router.POST("/users", func(c *gin.Context) {
		var reqBody userReqData
		c.BindJSON(&reqBody)
		c.JSON(http.StatusOK, resData{"user-id": reqBody.UserID})
	})

	router.PUT("/users", func(c *gin.Context) {
		var reqBody userReqData
		c.BindJSON(&reqBody)
		c.JSON(
			http.StatusOK,
			resData{"user-id": reqBody.UserID, "email": reqBody.Email},
		)
	})

	router.DELETE("/users", func(c *gin.Context) {
		var reqBody userReqData
		c.BindJSON(&reqBody)
		c.JSON(
			http.StatusOK,
			resData{"user-id": reqBody.UserID},
		)
	})
}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}
