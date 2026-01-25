package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/middlewares"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterCatalogProductRoutes(router *gin.Engine) {
	catalogProductService := services.NewCatalogProductService()
	userService := services.NewUserService()

	catalogProductController := controllers.NewCatalogProductController(
		catalogProductService,
		userService,
	)

	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		catalogs := protected.Group("/catalogs")
		{
			catalogs.GET("/:id/products", catalogProductController.ListCatalogProducts)
			catalogs.POST("/:id/attach-product", catalogProductController.AttachProduct)
			catalogs.PATCH("/:catalog_id/update-product/:id", catalogProductController.UpdateCatalogProduct)
			catalogs.DELETE("/:catalog_id/detach-product/:id", catalogProductController.DetachProduct)
		}
	}
}
