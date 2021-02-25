package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// GetAllPosts returns all posts
// optional: with tags or/and author
func GetAllPosts() gin.HandlerFunc {
	return func(c *gin.Context) {

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
}

// GetPost returns a post with given ID
func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, &response{
			"result": id,
		})
	}
}

// GetLikesForPost returns the total like count
// of a post with given ID
func GetLikesForPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = c.Param("id")
		c.JSON(200, &response{
			"result": 10,
		})
	}
}

// CreatePost creates a post in db
func CreatePost(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody postReqData
		c.BindJSON(&reqBody)
		c.JSON(http.StatusOK, &response{"id": reqBody.ID})
	}
}

// UpdatePost updates a post in db
func UpdatePost(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody postReqData
		c.BindJSON(&reqBody)
		c.JSON(
			http.StatusOK,
			&response{"id": reqBody.ID, "doc": reqBody.Doc},
		)
	}
}

// DeletePost deletes a post with given ID in db
func DeletePost(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody postReqData
		c.BindJSON(&reqBody)
		c.JSON(
			http.StatusOK,
			&response{"id": reqBody.ID},
		)
	}
}
