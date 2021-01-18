package posts

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type POSTdata struct {
	PostID string `json:"postid"`
}

// AddRoutes ...
// Adds HTTP Methods for the /posts endpoint
func AddRoutes(router *gin.Engine) {
	router.GET("/posts", func(c *gin.Context) {

		// Check queries
		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
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

	router.POST("/posts", func(c *gin.Context) {
		var postData POSTdata
		c.BindJSON(&postData)
		c.JSON(http.StatusOK, gin.H{"postid": postData.PostID})
	})
}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}
