package api

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

type postData struct {
	Doc      string `json:"doc" validate:"required"`
	Tags     string `json:"tags"`
	Likes    uint   `json:"likes"`
	Comments string `json:"comments"`
}

type newPostData struct {
	ID       uint   `json:"id" validate:"required"`
	Doc      string `json:"doc"`
	Tags     string `json:"tags"`
	Likes    uint   `json:"likes"`
	Comments string `json:"comments"`
}

type userUpdateForm struct {
	ID       uint   `json:"id" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type credential struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type errorResponse struct {
	Msg string `json:"message"`
}

type userUpdateQuery struct {
	ID             uint
	Email          string
	Password       string
	TokenExpiresIn interface{}
}

type postUpdateQuery struct {
	ID        uint
	Author    string
	Document  string
	Tags      string
	Comments  string
	Likes     uint
	CreatedAt time.Time
}

type response map[string]interface{}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}

func serializeUser(u *dbtool.User) response {
	return response{
		"id":    u.ID,
		"email": u.Email,
	}
}

func serializePost(p *dbtool.Post) response {
	return response{
		"id":       p.ID,
		"author":   p.Author,
		"doc":      p.Document,
		"tags":     p.Tags,
		"comments": p.Comments,
		"likes":    p.Likes,
	}
}

func createUserUpdateQuery(id uint, email, password string, tokenExpiresIn interface{}) (userUpdateQuery, error) {
	query := userUpdateQuery{
		ID: id,
	}

	if email == "" && password == "" {
		return query, errors.New("User update failed. No new data")
	}

	if email != "" {
		v := validator.New()

		if err := v.Var(email, "email"); err != nil {
			return query, errors.New("User update failed. Invalid email")
		}

		query.Email = email
	}

	if password != "" {
		query.Password = password
	}

	if tokenExpiresIn != nil {
		query.TokenExpiresIn = tokenExpiresIn
	}

	return query, nil
}

func validateStruct(c interface{}) error {
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}

// HandleError attaches error response to gin.Context
func HandleError(c *gin.Context, code int, msg string) {
	c.JSON(code, &errorResponse{Msg: msg})
}

func compareWithOldData(new, old interface{}, isString bool) (interface{}, bool) {
	changed := false
	if isString {
		if new != "" && new != old {
			changed = true
			return new, changed
		}
		return old, changed
	}

	if new != old {
		changed = true
		return new, changed
	}
	return old, changed

}

func createPostQuery(d *newPostData, currPost *dbtool.Post) (*postUpdateQuery, error) {
	doc, docChanged := compareWithOldData(d.Doc, currPost.Document, true)
	comments, comChanged := compareWithOldData(d.Comments, currPost.Comments, true)
	tags, tagsChanged := compareWithOldData(d.Tags, currPost.Tags, true)
	likes, likesChanged := compareWithOldData(d.Likes, currPost.Likes, false)
	query := postUpdateQuery{
		ID:        d.ID,
		Author:    currPost.Author,
		CreatedAt: currPost.CreatedAt,
		Document:  doc.(string),
		Comments:  comments.(string),
		Tags:      tags.(string),
		Likes:     likes.(uint),
	}

	if !docChanged && !comChanged && !tagsChanged && !likesChanged {
		return nil, fmt.Errorf("No new data to update")
	}
	return &query, nil
}

func checkIfUserIsAuthor(c *gin.Context, author string) bool {
	username, exists := c.Get("username")
	if !exists {
		return false

	}
	return username == author
}
