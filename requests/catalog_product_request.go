package requests

type AttachProductToCatalogRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Price     float64 `json:"price" binding:"required,min=0"`
	IsBought  bool    `json:"is_bought" binding:"omitempty"`
}

type MarkCatalogProductAsBoughtRequest struct {
	CatalogID uint `json:"catalog_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
}

type UpdateCatalogProductRequest struct {
	Price    *float64 `json:"price" binding:"omitempty,min=0"`
	IsBought *bool    `json:"is_bought" binding:"omitempty"`
}

type CreateExclusiveProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Images      []string `json:"images" binding:"required,min=1"`
	Price       float64  `json:"price" binding:"required,min=0"`
	IsBought    bool     `json:"is_bought" binding:"omitempty"`
}
