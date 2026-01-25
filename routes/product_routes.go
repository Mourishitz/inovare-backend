package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/middlewares"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.Engine) {
	productService := services.NewProductService()
	userService := services.NewUserService()
	productController := controllers.NewProductController(productService, userService)

	// All product routes require authentication and admin role (Role 2+)
	protected := router.Group("/api")
	protected.Use(middlewares.Authenticate())
	{
		products := protected.Group("/products")
		{
			products.GET("", productController.ListProducts)
			products.GET("/:id", productController.GetProduct)
			products.POST("", productController.CreateProduct)
			products.PATCH("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}
	}
}
