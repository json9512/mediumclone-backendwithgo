package users

import "github.com/gin-gonic/gin"

// UserSerializer holds reference to gin.Context
// for serializing data
type UserSerializer struct {
	context *gin.Context
}

// userResponse is format for sending user
// detail as JSON response to the client
type userResponse struct {
	ID    uint   `json:"user-id"`
	Email string `json:"email"`
}

// errorResponse is the uniform format for sending
// error to the client
type errorResponse struct {
	Msg string `json:"message"`
}

func serialize(u *User) userResponse {
	return userResponse{
		ID:    u.ID,
		Email: u.Email,
	}
}

// CreateUserData converts data given into
// User model for db
func CreateUserData(req userReqData, accessToken, refreshToken string) *User {
	return &User{
		ID:           req.UserID,
		Email:        req.Email,
		Password:     req.Password,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
