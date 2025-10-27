package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/middlewares"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	userService := services.NewUserService()
	authController := controllers.NewAuthController(userService)

	public := router.Group("/api")
	{
		public.POST("/login", authController.Login)
	}

	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		protected.GET("/me", authController.GetSessionUser)
	}
}
