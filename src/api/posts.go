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

// GetPosts godoc
// @Summary Get posts
// @Tags posts
// @Description Retrieve posts from db
// @ID get-posts
// @Accept  json
// @Produce  json
// @Param tags query string false "tags"
// @Param author query string false "author"
// @Success 200 {object} api.SwaggerPosts
// @Failure 400 {object} api.APIError "Bad Request"
// @Router /posts [get]
func GetPosts(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryChecker := &queryChecker{false, false}
		queries := c.Request.URL.Query()
		res := handleQuery(c, pool, queries, queryChecker)
		if !res {
			posts, err := db.GetPosts(c, pool)
			if err != nil {
				HandleError(c, http.StatusBadRequest, "No posts in db.")
			} else {
				c.JSON(200, serializePosts(*posts))
			}
		}
	}
}

// GetPost godoc
// @Summary Get post
// @Description Retrieve a post by its ID
// @Tags posts
// @ID get-post
// @Accept  json
// @Produce  json
// @Param id path string true "Post ID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} api.APIError "Bad Request"
// @Router /posts/{id} [get]
func GetPost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id := convertToInt(idStr)
		if id < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if post, err := db.GetPostByID(c, pool, id); err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found.")
		} else {
			c.JSON(http.StatusOK, serializePost(post))
		}
	}
}

// GetLikesForPost godoc
// @Summary Get likes of a post
// @Description Get like count of a post by its ID
// @Tags posts
// @ID get-post-likes
// @Accept  json
// @Produce  json
// @Param id path string true "Post ID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} api.APIError "Bad Request"
// @Router /posts/{id}/like [get]
func GetLikesForPost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id := convertToInt(idStr)
		if id < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		if likes, err := db.GetLikesForPost(c, pool, id); err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found.")
		} else {
			c.JSON(http.StatusOK, &response{"likes": likes})
		}
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new post in DB
// @Tags posts
// @ID create-post
// @Accept  json
// @Produce  json
// @Param post body api.PostInsertForm true "Add Post"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} api.APIError "Bad Request"
// @Failure 401 {object} api.APIError "Unauthorized"
// @Router /posts [post]
func CreatePost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody PostInsertForm
		if err := extractData(c, &reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid request data.")
			return
		}

		if err := validateStruct(&reqBody); err != nil || reqBody.Doc == "" {
			HandleError(c, http.StatusBadRequest, "ID, Doc required.")
			return
		}

		username, exists := c.Get("username")
		if !exists {
			HandleError(c, http.StatusBadRequest, "Username not found.")
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

// UpdatePost godoc
// @Summary Update a post
// @Description Updates a post by the provided data
// @Tags posts
// @ID update-post
// @Accept  json
// @Produce  json
// @Param post body api.PostUpdateForm true "Update Post"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} api.APIError "Bad Request"
// @Failure 401 {object} api.APIError "Unauthorized"
// @Router /posts [put]
func UpdatePost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody PostUpdateForm
		if err := extractData(c, &reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "Invalid request data.")
			return
		}

		if err := validateStruct(&reqBody); err != nil {
			HandleError(c, http.StatusBadRequest, "ID required.")
			return
		}
		postID := int64(reqBody.ID)
		queriedPost, err := db.GetPostByID(c, pool, postID)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found.")
			return
		}

		if !checkIfUserIsAuthor(c, queriedPost.Author.String) {
			HandleError(c, http.StatusBadRequest, "User is not the author of the post.")
			return
		}

		post, err := bindUpdateFormToPost(&reqBody, queriedPost.Author.String)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Update form not valid.")
			return
		}

		if createdPost, err := db.UpdatePost(c, pool, postID, post); err != nil {
			HandleError(c, http.StatusInternalServerError, "Failed to update post in DB.")
		} else {
			c.JSON(http.StatusOK, serializePost(createdPost))
		}
	}
}

// DeletePost godoc
// @Summary Delete a post
// @Description Deletes a post by its ID
// @Tags posts
// @ID delete-post
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} api.APIError "Bad Request"
// @Failure 401 {object} api.APIError "Unauthorized"
// @Router /posts/{id} [delete]
func DeletePost(pool *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id := convertToInt(idStr)
		if id < 1 {
			HandleError(c, http.StatusBadRequest, "Invalid ID.")
			return
		}

		queriedPost, err := db.GetPostByID(c, pool, id)
		if err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found.")
			return
		}

		if !checkIfUserIsAuthor(c, queriedPost.Author.String) {
			HandleError(c, http.StatusBadRequest, "User is not the author of the post.")
			return
		}

		if _, err := db.DeletePostByID(c, pool, id); err != nil {
			HandleError(c, http.StatusBadRequest, "Post not found.")
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
	rawAuthor, exists := q["author"]
	if !exists {
		return nil, false
	}

	if len(rawAuthor) > 0 {
		author := strings.ToLower(rawAuthor[0])
		return &author, true
	} else {
		author := strings.ToLower(rawAuthor[0])
		return &author, true
	}
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

func handleQuery(c *gin.Context, pool *sql.DB, q url.Values, checker *queryChecker) bool {
	success := false
	if checkIfQueriesExist(q) {
		tags, tagsExist := checkIfTagsExist(q)
		author, authorExists := checkIfAuthorExists(q)
		checker.authorExists = authorExists
		checker.tagsExist = tagsExist
		posts, err := executeQuery(c, pool, tags, author, checker)

		if err != nil || posts == nil {
			HandleError(c, http.StatusBadRequest, "No posts in db.")
		} else {
			c.JSON(200, serializePosts(*posts))
			success = true
		}
	}
	return success
}
