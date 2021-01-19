package posts

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type resData map[string]interface{}

type dataPOST struct {
	PostID string `json:"post-id"`
}

// AddRoutes adds HTTP Methods for the /posts endpoint
func AddRoutes(router *gin.Engine) {
	router.GET("/posts", func(c *gin.Context) {

		// Check queries
		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
			c.JSON(200, resData{
				"result": queries,
			})
		} else {
			c.JSON(200, resData{
				"result": []string{"test", "sample", "post"},
			})
		}
	})

	router.GET("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, resData{
			"result": id,
		})
	})

	router.GET("/posts/:id/like", func(c *gin.Context) {
		_ = c.Param("id")
		c.JSON(200, resData{
			"result": 10,
		})
	})

	router.POST("/posts", func(c *gin.Context) {
		var postData dataPOST
		c.BindJSON(&postData)
		c.JSON(http.StatusOK, resData{"post-id": postData.PostID})
	})
}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}
