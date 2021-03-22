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

	router.POST("/login", api.Login(db, env))
	router.POST("/logout", middlewares.VerifyUser(db), api.Logout(db))

	postsRouter := router.Group("/posts")
	postsRouter.GET("", api.GetAllPosts(db))
	postsRouter.GET("/:id", api.GetPost())
	postsRouter.GET("/:id/like", api.GetLikesForPost())
	postsRouter.POST("", middlewares.VerifyUser(db), api.CreatePost(db))
	postsRouter.PUT("", middlewares.VerifyUser(db), api.UpdatePost(db))
	postsRouter.DELETE("/:id", middlewares.VerifyUser(db), api.DeletePost(db))

	usersRouter := router.Group("/users")
	usersRouter.GET("/:id", api.RetrieveUser(db))
	usersRouter.POST("", api.RegisterUser(db))
	usersRouter.PUT("", middlewares.VerifyUser(db), api.UpdateUser(db))
	usersRouter.DELETE("/:id", middlewares.VerifyUser(db), api.DeleteUser(db))

}
