package requests

type AttachProductToCatalogRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Price     float64 `json:"price" binding:"required,min=0"`
	IsBought  bool    `json:"is_bought" binding:"omitempty"`
}

type UpdateCatalogProductRequest struct {
	Price    *float64 `json:"price" binding:"omitempty,min=0"`
	IsBought *bool    `json:"is_bought" binding:"omitempty"`
}
