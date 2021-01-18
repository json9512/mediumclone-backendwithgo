package users

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

// AddRoutes ...
// Adds HTTP Methods for the /users endpoint
func AddRoutes(router *gin.Engine) {
	router.GET("/users", func(c *gin.Context) {

		// Check queries
		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
			c.JSON(200, gin.H{
				"result": queries,
			})
		} else {
			c.JSON(200, gin.H{
				"result": []string{"test", "sample", "users"},
			})
		}
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, gin.H{
			"result": id,
		})
	})

}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}
