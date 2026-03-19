package controllers

import (
	"errors"
	"net/http"
	"strings"

	"inovare-backend/requests"
	"inovare-backend/services"
	"inovare-backend/utils"

	"github.com/gin-gonic/gin"
)

type PaymentWebhookController struct {
	catalogProductService services.CatalogProductService
}

func NewPaymentWebhookController(catalogProductService services.CatalogProductService) *PaymentWebhookController {
	return &PaymentWebhookController{
		catalogProductService: catalogProductService,
	}
}

// HandleCatalogPurchase handles POST /api/webhooks/payments.
func (c *PaymentWebhookController) HandleCatalogPurchase(ctx *gin.Context) {
	var req requests.PaymentWebhookRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	if req.Event != "billing.paid" || !strings.EqualFold(req.Data.Billing.Status, "PAID") {
		ctx.JSON(http.StatusOK, gin.H{"message": "Webhook event ignored"})
		return
	}

	externalIDs := make([]string, 0, len(req.Data.Billing.Products))
	for _, product := range req.Data.Billing.Products {
		externalIDs = append(externalIDs, product.ExternalID)
	}

	updatedProducts, err := c.catalogProductService.MarkAsBoughtByExternalIDs(externalIDs)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidWebhookExternalID):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product external ID"})
		case errors.Is(err, utils.ErrCatalogProductNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog product not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Webhook processed successfully",
		"products": updatedProducts,
	})
}
