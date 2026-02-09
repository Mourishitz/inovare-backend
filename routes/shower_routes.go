package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/middlewares"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterShowerRoutes(router *gin.Engine) {
	showerService := services.NewShowerService()
	userService := services.NewUserService()
	showerController := controllers.NewShowerController(showerService, userService)

	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		showers := protected.Group("/showers")
		{
			showers.GET("", showerController.ListShowers)
			showers.GET("/me", showerController.GetMyShowers)
			showers.GET("/:id", showerController.GetShower)
			showers.POST("", showerController.CreateShower)
			showers.PATCH("/:id", showerController.UpdateShower)
			showers.GET("/:id/catalog", showerController.GetShowerCatalog)
			showers.POST("/:id/catalog", showerController.AddCatalog)
			showers.POST("/:id/preferences", showerController.AddPreferences)
		}
	}
}
