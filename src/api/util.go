package api

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/models"
)

type postInsertForm struct {
	Doc      string `json:"doc" validate:"required"`
	Tags     string `json:"tags"`
	Likes    uint   `json:"likes"`
	Comments string `json:"comments"`
}

type postUpdateForm struct {
	ID       json.Number `json:"id" validate:"required"`
	Doc      string      `json:"doc"`
	Tags     string      `json:"tags"`
	Likes    uint        `json:"likes"`
	Comments string      `json:"comments"`
}

type userUpdateForm struct {
	ID             json.Number `json:"id" validate:"required"`
	Email          string      `json:"email"`
	Password       string      `json:"password"`
	TokenExpiresIn int64       `json:"token_expires_in"`
}

type userInsertForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type errorResponse struct {
	Msg string `json:"message"`
}

type response map[string]interface{}

// HandleError attaches error response to gin.Context
func HandleError(c *gin.Context, code int, msg string) {
	c.JSON(code, &errorResponse{Msg: msg})
}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}

func serializeUser(u *models.User) response {
	return response{
		"id":    u.ID,
		"email": u.Email,
	}
}

func serializePost(p *models.Post) response {
	author := strings.Title(strings.ToLower(p.Author.String))
	return response{
		"id":       p.ID,
		"author":   author,
		"doc":      p.Document,
		"tags":     p.Tags,
		"comments": p.Comments,
		"likes":    p.Likes,
	}
}

func serializePosts(posts []*models.Post) response {
	return response{
		"totalCount": len(posts),
		"posts":      posts,
	}
}

func extractData(c *gin.Context, reqBody interface{}) error {
	if err := c.BindJSON(&reqBody); err != nil {
		return err
	}
	return nil
}

func validateStruct(c interface{}) error {
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}

func checkIfUserIsAuthor(c *gin.Context, author string) bool {
	username, exists := c.Get("username")
	if !exists {
		return false
	}
	return username == author
}

func convertToInt(id string) int64 {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return -1
	}
	return idInt
}

func bindFormToPost(f *postInsertForm, author string) *db.Post {
	return &db.Post{
		Author:   strings.ToLower(author),
		Doc:      f.Doc,
		Comments: f.Comments,
		Tags:     strings.Split(f.Tags, ","),
		Likes:    int(f.Likes),
	}
}

func bindUpdateFormToPost(f *postUpdateForm, author string) (*db.Post, error) {
	var post db.Post
	postID, _ := f.ID.Int64()
	if postID < 0 {
		return nil, errors.New("ID required.")
	}

	if f.Comments == "" && f.Doc == "" && f.Likes < 0 && f.Tags == "" {
		return nil, errors.New("No new data.")
	}

	if f.Comments != "" {
		post.Comments = f.Comments
	}
	if f.Doc != "" {
		post.Comments = f.Comments
	}
	if f.Likes >= 0 {
		post.Likes = int(f.Likes)
	}
	if f.Tags != "" {
		post.Tags = strings.Split(f.Tags, ",")
	}
	post.Author = author

	return &post, nil
}

func bindFormToUser(f *userInsertForm) *db.User {
	return &db.User{
		Email:          f.Email,
		Password:       f.Password,
		TokenExpiresIn: 0,
	}

}

func bindUpdateFormToUser(b *userUpdateForm) (*db.User, error) {
	var user db.User
	userID, _ := b.ID.Int64()
	if userID < 0 {
		return nil, errors.New("ID required.")
	}

	if b.Email == "" && b.Password == "" && b.TokenExpiresIn < 0 {
		return nil, errors.New("No new data.")
	}

	if b.Email != "" {
		v := validator.New()

		if err := v.Var(b.Email, "email"); err != nil {
			return nil, errors.New("Invalid email.")
		}

		user.Email = b.Email
	}

	if b.Password != "" {
		user.Password = b.Password
	}

	if b.TokenExpiresIn > -1 {
		user.TokenExpiresIn = b.TokenExpiresIn
	}

	return &user, nil
}
