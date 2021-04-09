package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/api"
	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/middlewares"
)

// AddRoutes adds available routes to the provided router
func AddRoutes(router *gin.Engine, db *sql.DB, env *config.EnvVars) {
	apiGroup := router.Group("/api/v1")
	{
		apiGroup.POST("/login", api.Login(db, env))
		apiGroup.POST("/logout", middlewares.VerifyUser(db), api.Logout(db))

		posts := apiGroup.Group("/posts")
		posts.GET("", api.GetPosts(db))
		posts.GET(":id", api.GetPost(db))
		posts.GET(":id/like", api.GetLikesForPost(db))
		posts.POST("", middlewares.VerifyUser(db), api.CreatePost(db))
		posts.PUT("", middlewares.VerifyUser(db), api.UpdatePost(db))
		posts.DELETE(":id", middlewares.VerifyUser(db), api.DeletePost(db))

		users := apiGroup.Group("/users")
		users.GET(":id", api.RetrieveUser(db))
		users.POST("", api.RegisterUser(db))
		users.PUT("", middlewares.VerifyUser(db), api.UpdateUser(db))
		users.DELETE(":id", middlewares.VerifyUser(db), api.DeleteUser(db))
	}
}
