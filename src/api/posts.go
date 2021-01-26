package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllPosts returns all posts
// optional: with tags or/and author
func GetAllPosts() gin.HandlerFunc {
	handler := func(c *gin.Context) {

		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
			c.JSON(200, &response{
				"result": queries,
			})
		} else {
			c.JSON(200, &response{
				"result": []string{"test", "sample", "post"},
			})
		}
	}
	return gin.HandlerFunc(handler)
}

// GetSinglePost returns a post with given ID
func GetSinglePost() gin.HandlerFunc {
	handler := func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, &response{
			"result": id,
		})
	}
	return gin.HandlerFunc(handler)
}

// GetLikesForPost returns the total like count
// of a post with given ID
func GetLikesForPost() gin.HandlerFunc {
	handler := func(c *gin.Context) {
		_ = c.Param("id")
		c.JSON(200, &response{
			"result": 10,
		})
	}
	return gin.HandlerFunc(handler)
}

// CreatePost creates a post in db
func CreatePost() gin.HandlerFunc {
	handler := func(c *gin.Context) {
		var reqBody PostReqData
		c.BindJSON(&reqBody)
		c.JSON(http.StatusOK, &response{"post-id": reqBody.PostID})
	}
	return gin.HandlerFunc(handler)
}

// UpdatePost updates a post in db
func UpdatePost() gin.HandlerFunc {
	handler := func(c *gin.Context) {
		var reqBody PostReqData
		c.BindJSON(&reqBody)
		c.JSON(
			http.StatusOK,
			&response{"post-id": reqBody.PostID, "doc": reqBody.Doc},
		)
	}
	return gin.HandlerFunc(handler)
}

// DeletePost deletes a post with given ID in db
func DeletePost() gin.HandlerFunc {
	handler := func(c *gin.Context) {
		var reqBody PostReqData
		c.BindJSON(&reqBody)
		c.JSON(
			http.StatusOK,
			&response{"post-id": reqBody.PostID},
		)
	}
	return gin.HandlerFunc(handler)
}
