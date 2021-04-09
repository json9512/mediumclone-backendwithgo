package api

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/models"
)

type PostInsertForm struct {
	Doc      string `json:"doc" validate:"required" example:"some-text"`
	Tags     string `json:"tags" example:"some,tags,here"`
	Likes    uint   `json:"likes" example:"123"`
	Comments string `json:"comments" example:"some-comment"`
}

type PostUpdateForm struct {
	ID       int    `json:"id" validate:"required" example:"1"`
	Doc      string `json:"doc" example:"some-text"`
	Tags     string `json:"tags" example:"some,tags,here"`
	Likes    uint   `json:"likes" example:"123"`
	Comments string `json:"comments" example:"some-comment"`
}

type UserUpdateForm struct {
	ID             int    `json:"id" example:"1" validate:"required"`
	Email          string `json:"email" example:"someone@somewhere.com"`
	Password       string `json:"password" example:"very-hard-password!2"`
	TokenExpiresIn int64  `json:"token_expires_in" example:"15233324"`
}

type UserInsertForm struct {
	Email    string `json:"email" example:"someone@somewhere.com" validate:"required,email"`
	Password string `json:"password" example:"very-hard-password!2" validate:"required"`
}

type APIError struct {
	Msg string `json:"message"`
}

type SwaggerPosts struct {
	TotalCount int              `json:"total_count"`
	Posts      []PostUpdateForm `json:"posts"`
}

type SwaggerUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type SwaggerEmail struct {
	Email string `json:"email" example:"someone@somewhere.com"`
}

type response map[string]interface{}

// HandleError attaches error response to gin.Context
func HandleError(c *gin.Context, code int, msg string) {
	c.JSON(code, &APIError{Msg: msg})
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
		"total_count": len(posts),
		"posts":       posts,
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

func bindFormToPost(f *PostInsertForm, author string) *db.Post {
	return &db.Post{
		Author:   strings.ToLower(author),
		Doc:      f.Doc,
		Comments: f.Comments,
		Tags:     strings.Split(f.Tags, ","),
		Likes:    int(f.Likes),
	}
}

func bindUpdateFormToPost(f *PostUpdateForm, author string) (*db.Post, error) {
	var post db.Post
	postID := int64(f.ID)
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

func bindFormToUser(f *UserInsertForm) *db.User {
	return &db.User{
		Email:          f.Email,
		Password:       f.Password,
		TokenExpiresIn: 0,
	}

}

func bindUpdateFormToUser(b *UserUpdateForm) (*db.User, error) {
	var user db.User
	userID := int64(b.ID)
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
