package requests

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url" binding:"required,url"`
}

type UpdateProductRequest struct {
	Name        *string `json:"name" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
	ImageURL    *string `json:"image_url" binding:"omitempty,url"`
}
