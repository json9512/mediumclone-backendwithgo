package users

import "github.com/gin-gonic/gin"

// UserSerializer holds reference to gin.Context
// for serializing data
type UserSerializer struct {
	context *gin.Context
}

// UserResponse is format for sending user
// detail as JSON response to the client
type UserResponse struct {
	ID           uint   `json:"user-id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ErrorResponse is the uniform format for sending
// error to the client
type ErrorResponse struct {
	Msg string `json:"message"`
}

// Serialize converts User model (db) to
// UserResponse format
func Serialize(data *User) UserResponse {
	return UserResponse{
		ID:           data.ID,
		Email:        data.Email,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}
}

// CreateUserData converts data given into
// User model for db
func CreateUserData(data userReqData, accessToken, refreshToken string) *User {
	return &User{
		ID:           data.UserID,
		Email:        data.Email,
		Password:     data.Password,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
