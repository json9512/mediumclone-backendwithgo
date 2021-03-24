package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// GetPosts returns all posts
// optional: with tags or/and author
func GetPosts(pool *sql.DB) gin.HandlerFunc {
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
				// 	query, _ := db.GetPostsByTags(tags)
				// 	fmt.Println(query)
				//
			}
			posts, err := db.GetPosts(c, pool)
			if err != nil {
				HandleError(c, http.StatusBadRequest, "No posts in db")
			}

			c.JSON(200, &response{
				"result": posts,
			})
		} else {
			c.JSON(200, &response{
				"result": []string{"test", "sample", "post"},
			})
		}
	}
}

// GetPost returns a post with given ID
func GetPost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id := convertToInt(idStr)
		if id < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if post, err := db.GetPostByID(c, pool, id); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid Request.")
			return
		} else {
			c.JSON(http.StatusOK, &response{
				"result": post, // need serialize
			})
		}
	}
}

// GetLikesForPost returns the total like count
// of a post with given ID
func GetLikesForPost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id := convertToInt(idStr)
		if id < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if likes, err := db.GetLikesForPost(c, pool, id); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid Request.")
		} else {
			c.JSON(http.StatusOK, likes) // need serialize of some sort
		}
	}
}

// CreatePost creates a post in db
func CreatePost(pool *sql.DB) gin.HandlerFunc {
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

		post := bindFormToPost(&reqBody, username.(string))
		if createdPost, err := db.InsertPost(c, pool, post); err != nil {
			HandleError(c, http.StatusInternalServerError, "Failed to create post in DB.")
		} else {
			c.JSON(http.StatusOK, serializePost(createdPost))
		}
	}
}

// UpdatePost updates a post in db
func UpdatePost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody postUpdateForm
		if err := extractData(c, &reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid request data")
			return
		}

		if err := validateStruct(&reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "ID required")
			return
		}

		queriedPost, err := db.GetPostByID(c, pool, reqBody.id)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found")
			return
		}

		if !checkIfUserIsAuthor(c, queriedPost.Author) {
			HandleError(c, http.StatusBadRequest, "User is not the author of the post")
			return
		}

		post := bindFormToPost(reqBody, queriedPost.Author)
		if createdPost, err := db.UpdatePost(c, pool, reqBody.id, post); err != nil {
			HandleError(c, http.StatusInternalServerError, "Failed to update post in DB.")
			return
		} else {
			c.JSON(http.StatusOK, serializePost(createdPost))
		}
	}
}

// DeletePost deletes a post with given ID in db
func DeletePost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id := convertToInt(idStr)
		if id < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if _, err := db.DeletePostByID(c, pool, id); err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found")
		} else {
			c.Status(http.StatusOK)
		}
	}
}
