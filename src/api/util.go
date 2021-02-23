package api

import (
	"errors"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

type postReqData struct {
	ID       string `json:"id"`
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

type userResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type updateQuery struct {
	ID             uint
	Email          string
	Password       string
	TokenExpiresIn interface{}
}

type response map[string]interface{}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}

func serializeUser(u *dbtool.User) userResponse {
	return userResponse{
		ID:    u.ID,
		Email: u.Email,
	}
}

func createUpdateQuery(id, email, password, tokenExpiresIn interface{}) (updateQuery, error) {
	query := updateQuery{
		ID: id.(uint),
	}

	if email == "" && password == "" {
		return query, errors.New("User update failed. No new data")
	}

	if email != "" {
		v := validator.New()

		if err := v.Var(email, "email"); err != nil {
			return query, errors.New("User update failed. Invalid email")
		}

		query.Email = email.(string)
	}

	if password != "" {
		query.Password = password.(string)
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
