package requests

type PaymentWebhookRequest struct {
	Event   string              `json:"event" binding:"required"`
	Data    *PaymentWebhookData `json:"data" binding:"required"`
	DevMode bool                `json:"devMode"`
}

type PaymentWebhookData struct {
	Billing *PaymentWebhookBilling `json:"billing" binding:"required"`
}

type PaymentWebhookBilling struct {
	ID       string                         `json:"id"`
	Status   string                         `json:"status" binding:"required"`
	Products []PaymentWebhookBillingProduct `json:"products" binding:"required,min=1,dive"`
}

type PaymentWebhookBillingProduct struct {
	PublicID   string `json:"publicId"`
	ExternalID string `json:"externalId" binding:"required"`
	Quantity   int    `json:"quantity"`
}
