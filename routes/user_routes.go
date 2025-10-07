package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/middlewares"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userService := services.NewUserService()
	userController := controllers.NewUserController(userService)

	public := router.Group("/api")
	{
		public.POST("/register", userController.CreateUser)

		users := public.Group("/users")
		{
			users.GET("/:id", userController.GetUser)
		}
	}

	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		users := protected.Group("/users")
		{
			users.POST("/", userController.CreateUser)
		}
	}
}
