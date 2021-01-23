package posts

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
)

type postReqData struct {
	PostID string `json:"post-id"`
	Doc    string `json:"doc"`
}

// AddRoutes adds HTTP Methods for the /posts endpoint
func AddRoutes(router *gin.Engine) {
	router.GET("/posts", func(c *gin.Context) {

		// Check queries
		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
			c.JSON(200, &config.ResData{
				"result": queries,
			})
		} else {
			c.JSON(200, &config.ResData{
				"result": []string{"test", "sample", "post"},
			})
		}
	})

	router.GET("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, &config.ResData{
			"result": id,
		})
	})

	router.GET("/posts/:id/like", func(c *gin.Context) {
		_ = c.Param("id")
		c.JSON(200, &config.ResData{
			"result": 10,
		})
	})

	router.POST("/posts", func(c *gin.Context) {
		var reqBody postReqData
		c.BindJSON(&reqBody)
		c.JSON(http.StatusOK, &config.ResData{"post-id": reqBody.PostID})
	})

	router.PUT("/posts", func(c *gin.Context) {
		var reqBody postReqData
		c.BindJSON(&reqBody)
		c.JSON(
			http.StatusOK,
			&config.ResData{"post-id": reqBody.PostID, "doc": reqBody.Doc},
		)
	})

	router.DELETE("/posts", func(c *gin.Context) {
		var reqBody postReqData
		c.BindJSON(&reqBody)
		c.JSON(
			http.StatusOK,
			&config.ResData{"post-id": reqBody.PostID},
		)
	})
}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}
