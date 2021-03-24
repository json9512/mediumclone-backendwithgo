package api

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/volatiletech/null"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/models"
)

type postForm struct {
	Doc      string `json:"doc" validate:"required"`
	Tags     string `json:"tags"`
	Likes    uint   `json:"likes"`
	Comments string `json:"comments"`
}

type postUpdateForm struct {
	id int64 `json:"id" validate:"required"`
	*postForm
}

type userUpdateForm struct {
	ID             uint   `json:"id" validate:"required"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	TokenExpiresIn int64  `json:"token_expires_in"`
}

type credential struct {
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
	return response{
		"id":       p.ID,
		"author":   p.Author,
		"doc":      p.Document,
		"tags":     p.Tags,
		"comments": p.Comments,
		"likes":    p.Likes,
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

func bindFormToPost(f interface{}, author string) *db.Post {
	if postUpdateForm, ok := f.(postUpdateForm); ok {
		return &db.Post{
			Author:   author,
			Doc:      postUpdateForm.Doc,
			Comments: postUpdateForm.Comments,
			Tags:     strings.Split(postUpdateForm.Tags, ","),
			Likes:    int(postUpdateForm.Likes),
		}
	}

	if postForm, ok := f.(postForm); ok {
		return &db.Post{
			Author:   author,
			Doc:      postForm.Doc,
			Comments: postForm.Comments,
			Tags:     strings.Split(postForm.Tags, ","),
			Likes:    int(postForm.Likes),
		}
	}
	return nil
}

func bindFormToUser(f interface{}) *db.User {
	if userForm, ok := f.(credential); ok {
		return &db.User{
			Email:          userForm.Email,
			Password:       userForm.Password,
			TokenExpiresIn: 0,
		}
	}
	return nil
}

func bindUpdateFormToUser(b *userUpdateForm) (*db.User, error) {
	var user db.User
	if b.ID < 0 {
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

		user.Email = null.StringFrom(b.Email)
	}

	if b.Password != "" {
		user.Password = null.StringFrom(b.Password)
	}

	if b.TokenExpiresIn > -1 {
		user.TokenExpiresIn = null.Int64From(b.TokenExpiresIn)
	}

	return &user, nil
}
