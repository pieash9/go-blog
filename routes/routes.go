package routes

import (
	"go-blog/controllers"
	"go-blog/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// auth
	auth := api.Group("/auth")
	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// posts route
	protected.POST("/posts", controllers.CreatePost)
	protected.GET("/posts", controllers.GetPosts)
	protected.GET("/posts/:id", controllers.GetPostById)
	protected.PUT("/posts/:id", controllers.UpdatePost)
	protected.DELETE("/posts/:id", controllers.DeletePost)

	// comments route
	protected.POST("/comments/:postId", controllers.CreateComment)
	protected.GET("/comments/:postId", controllers.GetComments)
	protected.DELETE("/comments/:commentId", controllers.DeleteComment)
}
