package posts

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

// AddRoutes ...
// Adds HTTP Methods for the /posts endpoint
func AddRoutes(router *gin.Engine) {
	router.GET("/posts", func(c *gin.Context) {

		// Check queries
		queries := c.Request.URL.Query()

		if checkIfExist(queries) {
			c.JSON(200, gin.H{
				"result": queries,
			})
		} else {
			c.JSON(200, gin.H{
				"result": []string{"test", "sample", "post"},
			})
		}
	})

	router.GET("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, gin.H{
			"result": id,
		})
	})

	router.GET("/posts/:id/like", func(c *gin.Context) {
		_ = c.Param("id")
		c.JSON(200, gin.H{
			"result": 10,
		})
	})
}

func checkIfExist(values url.Values) bool {
	if len(values) > 0 {
		return true
	}
	return false
}
