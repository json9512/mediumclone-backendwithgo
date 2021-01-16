package posts

import (
	"github.com/gin-gonic/gin"
)

// PostsRouter ...
// Adds HTTP Methods for the /posts endpoint
func PostsRouter(router *gin.Engine) {
	router.GET("/posts", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": []string{"test", "sample", "post"},
		})
	})

	router.GET("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, gin.H{
			"result": id,
		})
	})
}
