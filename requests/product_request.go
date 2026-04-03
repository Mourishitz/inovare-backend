package requests

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Images      []string `json:"images" binding:"required,min=1"`
	IsExclusive bool     `json:"is_exclusive" binding:"omitempty"`
	CatalogID   *uint    `json:"catalog_id" binding:"omitempty"`
}

type UpdateProductRequest struct {
	Name        *string   `json:"name" binding:"omitempty"`
	Description *string   `json:"description" binding:"omitempty"`
	Images      *[]string `json:"images" binding:"omitempty"`
	IsExclusive *bool     `json:"is_exclusive" binding:"omitempty"`
}
