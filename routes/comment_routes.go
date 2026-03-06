package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/middlewares"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(router *gin.Engine) {
	commentService := services.NewCommentService()
	commentController := controllers.NewCommentController(commentService)

	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		catalogs := protected.Group("/catalogs")
		{
			catalogs.POST("/:id/comments", commentController.AddComment)
			catalogs.GET("/:id/comments", commentController.ListComments)
		}
	}
}
