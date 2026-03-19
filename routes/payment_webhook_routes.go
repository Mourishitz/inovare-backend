package routes

import (
	"inovare-backend/controllers"
	"inovare-backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentWebhookRoutes(router *gin.Engine) {
	catalogProductService := services.NewCatalogProductService()
	paymentWebhookController := controllers.NewPaymentWebhookController(catalogProductService)

	public := router.Group("/api/webhooks")
	{
		public.POST("/payments", paymentWebhookController.HandleCatalogPurchase)
	}
}
