package requests

type PublicCatalogPurchaseRequest struct {
	ProductID  string                        `json:"productId" binding:"required"`
	Customer   PublicCatalogPurchaseCustomer `json:"customer" binding:"required"`
	CurrentURL string                        `json:"currentUrl" binding:"required"`
}

type PublicCatalogPurchaseCustomer struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Cellphone string `json:"cellphone" binding:"required"`
	TaxID     string `json:"taxId" binding:"required"`
}

type PublicCatalogPurchaseResponse struct {
	CheckoutURL string `json:"checkoutUrl"`
}
