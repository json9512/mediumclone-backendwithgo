package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/api"
	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
	"github.com/json9512/mediumclone-backendwithgo/src/middlewares"
)

// AddRoutes adds available routes to the provided router
func AddRoutes(router *gin.Engine, db *dbtool.DB, envVars *config.EnvVars) {

	router.POST("/login", api.Login(db, envVars))
	router.POST("/logout", middlewares.VerifyUser(db), api.Logout(db))

	postsRouter := router.Group("/posts")
	postsRouter.GET("", api.GetAllPosts())
	postsRouter.GET("/:id", api.GetPost())
	postsRouter.GET("/:id/like", api.GetLikesForPost())
	postsRouter.POST("", middlewares.VerifyUser(db), api.CreatePost(db))
	postsRouter.PUT("", api.UpdatePost())
	postsRouter.DELETE("", api.DeletePost())

	usersRouter := router.Group("/users")
	usersRouter.GET("/:id", api.RetrieveUser(db))
	usersRouter.POST("", api.RegisterUser(db))
	usersRouter.PUT("", middlewares.VerifyUser(db), api.UpdateUser(db))
	usersRouter.DELETE("/:id", middlewares.VerifyUser(db), api.DeleteUser(db))

}
