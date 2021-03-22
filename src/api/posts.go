package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// GetAllPosts returns all posts
// optional: with tags or/and author
func GetAllPosts(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
			// extract tags
			// fmt.Println(queries)
			// 3 condition
			// 1. only tags exist
			tags, tagsExist := queries["tags"]
			if tagsExist && strings.Contains(tags[0], ",") {
				tags = strings.Split(tags[0], ",")
			}
			// 2. only author exists (atmost 1)
			var author string
			rawAuthor, _ := queries["author"]
			if len(rawAuthor) > 0 {
				author = rawAuthor[0]
			}
			// 3. both exist
			if len(tags) > 0 && author != "" {
				query, _ := db.GetPostsByTags(tags)
				fmt.Println(query)
			}

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
		var reqBody postForm
		if err := extractData(c, &reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid request data")
			return
		}

		if err := validateStruct(&reqBody); err != nil && reqBody.Doc == "" {
			HandleError(c, http.StatusBadRequest, "ID, Doc required")
			return
		}

		username, exists := c.Get("username")
		if !exists {
			HandleError(c, http.StatusBadRequest, "Username not found")
			return
		}

		post, err := db.CreatePost(reqBody.Doc, username.(string), strings.Split(reqBody.Tags, ","))
		if err != nil {
			HandleError(c, http.StatusInternalServerError, "Failed to create post in DB.")
			return
		}

		c.JSON(http.StatusOK, serializePost(post))
	}
}

// UpdatePost updates a post in db
func UpdatePost(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody dbtool.UpdatePostForm
		if err := extractData(c, &reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid request data")
			return
		}

		if err := validateStruct(&reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "ID required")
			return
		}

		post, err := db.GetPostByID(reqBody.ID)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found")
			return
		}

		if !checkIfUserIsAuthor(c, post.Author) {
			HandleError(c, http.StatusBadRequest, "User is not the author of the post")
			return
		}

		if _, err := db.UpdatePost(post.ID, reqBody); err != nil {
			HandleError(c, http.StatusInternalServerError, "Failed to update post in DB.")
			return
		}

		// Check that the post's author is this user
		c.Status(http.StatusOK)
	}
}

// DeletePost deletes a post with given ID in db
func DeletePost(db *dbtool.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)

		if err != nil || idInt < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID")
			return
		}

		if _, err := db.DeletePostByID(idInt); err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found")
		} else {
			c.Status(http.StatusOK)
		}
	}
}
