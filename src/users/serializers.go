package users

import "github.com/gin-gonic/gin"

type UserSerializer struct {
	context *gin.Context
}

type UserResponse struct {
	ID           uint   `json:"user-id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	msg string `json:"message"`
}

func Serialize(data *User) UserResponse {
	return UserResponse{
		ID:           data.ID,
		Email:        data.Email,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}
}

func CreateUserData(data userReqData, accessToken, refreshToken string) *User {
	// Token validation or generation needed here?
	return &User{
		ID:           data.UserID,
		Email:        data.Email,
		Password:     data.Password,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
