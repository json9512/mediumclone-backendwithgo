package api

import (
	"errors"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

type postReqData struct {
	ID  string `json:"id"`
	Doc string `json:"doc"`
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

type CustomError struct {
	G    *gin.Context
	Code int
	Msg  string
}

type errorResponse struct {
	Msg string `json:"message"`
}

type userResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
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

func createRegUser(cred credential) dbtool.User {
	return dbtool.User{
		ID:       0,
		Email:    cred.Email,
		Password: cred.Password,
	}
}

func createUserWithNewData(u userUpdateForm) (dbtool.User, error) {
	user := dbtool.User{
		ID: u.ID,
	}

	if u.Email == "" && u.Password == "" {
		return user, errors.New("User update failed. No new data")
	}

	if u.Email != "" {
		v := validator.New()

		if err := v.Var(u.Email, "email"); err != nil {
			return user, errors.New("User update failed. Invalid email")
		}

		user.Email = u.Email
	}

	if u.Password != "" {
		user.Password = u.Password
	}

	return user, nil
}

func validateStruct(c interface{}) error {
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}

func HandleError(c *gin.Context, code int, msg string) {
	c.JSON(code, &errorResponse{Msg: msg})
}
