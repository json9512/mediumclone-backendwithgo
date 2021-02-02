package api

import (
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

type PostReqData struct {
	PostID string `json:"post-id"`
	Doc    string `json:"doc"`
}

// 이게 필요한가 ?
type UserReqData struct {
	UserID   uint   `json:"user-id"`
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
	ID    uint   `json:"user-id"`
	Email string `json:"email"`
}

type response map[string]interface{}

func checkIfQueriesExist(v url.Values) bool {
	if len(v) > 0 {
		return true
	}
	return false
}

func serializeUser(u dbtool.User) userResponse {
	return userResponse{
		ID:    u.ID,
		Email: u.Email,
	}
}

func createUserObj(req UserReqData) dbtool.User {
	return dbtool.User{
		ID:       req.UserID,
		Email:    req.Email,
		Password: req.Password,
	}
}

func validateCredential(c *credential) error {
	v := validator.New()
	if valErr := v.Struct(c); valErr != nil {
		return valErr
	}
	return nil
}
