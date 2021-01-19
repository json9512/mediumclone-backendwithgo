package users

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type resData map[string]interface{}

type dataPOST struct {
	UserID string `json:"user-id"`
}

// AddRoutes adds HTTP Methods for the /users endpoint
func AddRoutes(router *gin.Engine) {
	router.GET("/users", func(c *gin.Context) {

		// Check queries
		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
			c.JSON(200, gin.H{
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
		var userData dataPOST
		c.BindJSON(&userData)
		c.JSON(http.StatusOK, resData{"user-id": userData.UserID})
	})
}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}
