package api

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/models"
)

type queryChecker struct {
	tagsExist    bool
	authorExists bool
}

// GetPosts returns all posts
// optional: with tags or/and author
func GetPosts(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryChecker := &queryChecker{false, false}
		queries := c.Request.URL.Query()

		if checkIfQueriesExist(queries) {
			tags, tagsExist := checkIfTagsExist(queries)
			author, authorExists := checkIfAuthorExists(queries)
			queryChecker.authorExists = authorExists
			queryChecker.tagsExist = tagsExist
			posts, err := executeQuery(c, pool, tags, author, queryChecker)

			if err != nil || posts == nil {
				HandleError(c, http.StatusBadRequest, "No posts in db")
			} else {
				c.JSON(200, serializePosts(*posts))
			}
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
			c.JSON(http.StatusOK, serializePost(post))
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
			c.JSON(http.StatusOK, &response{"likes": likes})
		}
	}
}

// CreatePost creates a post in db
func CreatePost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody postInsertForm
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
		postID, _ := reqBody.ID.Int64()
		queriedPost, err := db.GetPostByID(c, pool, postID)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found")
			return
		}

		if !checkIfUserIsAuthor(c, queriedPost.Author.String) {
			HandleError(c, http.StatusBadRequest, "User is not the author of the post")
			return
		}

		post, err := bindUpdateFormToPost(&reqBody, queriedPost.Author.String)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Update form not valid.")
			return
		}

		if createdPost, err := db.UpdatePost(c, pool, postID, post); err != nil {
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

func checkIfTagsExist(q url.Values) (*[]string, bool) {
	tags, tagsExist := q["tags"]
	if !tagsExist {
		return nil, false
	}

	if tagsExist && strings.Contains(tags[0], ",") {
		tags = strings.Split(tags[0], ",")
		return &tags, true
	} else {
		return &tags, true
	}
}

func checkIfAuthorExists(q url.Values) (*string, bool) {
	var author string
	rawAuthor, _ := q["author"]
	if len(rawAuthor) > 0 {
		author = rawAuthor[0]
		return &author, true
	}
	return nil, false
}

func executeQuery(c context.Context, pool *sql.DB, tags *[]string, author *string, checker *queryChecker) (*models.PostSlice, error) {
	if checker.authorExists && checker.tagsExist {
		return db.GetPostsByTagsAndFilterByAuthor(c, pool, *tags, *author)
	} else if checker.authorExists {
		return db.GetPostsByAuthor(c, pool, *author)
	} else if checker.tagsExist {
		return db.GetPostsByTags(c, pool, *tags)
	}
	return nil, nil
}
