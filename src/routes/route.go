package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/api"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// AddRoutes adds available routes to the provided router
func AddRoutes(router *gin.Engine, db *dbtool.Pool) {
	router.POST("/login", api.Login(db))
	router.POST("/logout", api.Logout(db))

	router.GET("/posts", api.GetAllPosts())
	router.GET("/posts/:id", api.GetSinglePost())
	router.GET("/posts/:id/like", api.GetLikesForPost())
	router.POST("/posts", api.CreatePost())
	router.PUT("/posts", api.UpdatePost())
	router.DELETE("/posts", api.DeletePost())

	router.GET("/users/:id", api.RetrieveUser(db))
	router.POST("/users", api.RegisterUser(db))
	router.PUT("/users", api.UpdateUser(db))
	router.DELETE("/users/:id", api.DeleteUser(db))
}
