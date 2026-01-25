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

	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		catalogs := protected.Group("/catalogs")
		{
			catalogs.PATCH("/:id/approve", catalogController.ApproveCatalog)
		}
	}
}
