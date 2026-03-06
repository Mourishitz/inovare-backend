package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/middlewares"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.Engine) {
	productService := services.NewProductService()
	catalogProductService := services.NewCatalogProductService()
	userService := services.NewUserService()
	productController := controllers.NewProductController(productService, catalogProductService, userService)

	// All product routes require authentication and admin role (Role 2+)
	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		products := protected.Group("/products")
		{
			products.GET("", productController.ListProducts)
			products.GET("/search", productController.SearchProducts)
			products.GET("/:id", productController.GetProduct)
			products.GET("/:id/image", productController.GetProductImage)
			products.POST("", productController.CreateProduct)
			products.PATCH("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}
	}
}
