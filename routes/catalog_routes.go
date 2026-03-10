package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/middlewares"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterCatalogRoutes(router *gin.Engine) {
	catalogService := services.NewCatalogService()
	userService := services.NewUserService()

	catalogController := controllers.NewCatalogController(catalogService, userService)

	// Public routes (no auth required)
	public := router.Group("/api/catalogs")
	{
		public.GET("/url/:url", catalogController.GetByURL)
	}

	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		catalogs := protected.Group("/catalogs")
		{
			catalogs.GET("/:id", catalogController.GetByID)
			catalogs.PATCH("/:id/approve", catalogController.ApproveCatalog)
			catalogs.PATCH("/:id/changes-made", catalogController.RegisterChanges)
		}
	}
}
