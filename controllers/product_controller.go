package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"inovare-backend/requests"
	"inovare-backend/services"
	"inovare-backend/utils"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService        services.ProductService
	catalogProductService services.CatalogProductService
	userService           services.UserService
}

func NewProductController(productService services.ProductService, catalogProductService services.CatalogProductService, userService services.UserService) *ProductController {
	return &ProductController{
		productService:        productService,
		catalogProductService: catalogProductService,
		userService:           userService,
	}
}

// checkAdminRole validates if the user is an admin (Role 2+)
func (c *ProductController) checkAdminRole(ctx *gin.Context) bool {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return false
	}

	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}

	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return false
	}

	return true
}

// ListProducts handles GET /api/products
func (c *ProductController) ListProducts(ctx *gin.Context) {
	if !c.checkAdminRole(ctx) {
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total, err := c.productService.GetAllPaginated(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type productListItem struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IsExclusive bool   `json:"is_exclusive"`
	}

	items := make([]productListItem, len(products))
	for i, p := range products {
		items[i] = productListItem{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			IsExclusive: p.IsExclusive,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": items,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// SearchProducts handles GET /api/products/search
func (c *ProductController) SearchProducts(ctx *gin.Context) {
	query := ctx.Query("q")

	var catalogID *uint
	if rawID := ctx.Query("catalog_id"); rawID != "" {
		id, err := strconv.ParseUint(rawID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog_id"})
			return
		}
		uid := uint(id)
		catalogID = &uid
	}

	products, err := c.productService.Search(query, catalogID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type searchResultItem struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
		IsExclusive bool   `json:"is_exclusive"`
	}

	items := make([]searchResultItem, len(products))
	for i, p := range products {
		items[i] = searchResultItem{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			ImageURL:    p.ImageURL,
			IsExclusive: p.IsExclusive,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

// GetProduct handles GET /api/products/:id
func (c *ProductController) GetProduct(ctx *gin.Context) {
	if !c.checkAdminRole(ctx) {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := c.productService.GetByID(id)
	if err != nil {
		if errors.Is(err, utils.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var catalogID *uint
	if product.IsExclusive {
		catalogID, err = c.catalogProductService.GetCatalogIDByProductID(product.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           product.ID,
		"name":         product.Name,
		"description":  product.Description,
		"image_url":    product.ImageURL,
		"is_exclusive": product.IsExclusive,
		"catalog_id":   catalogID,
	})
}

// GetProductImage handles GET /api/products/:id/image
func (c *ProductController) GetProductImage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := c.productService.GetByID(id)
	if err != nil {
		if errors.Is(err, utils.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"image_url": product.ImageURL})
}

// CreateProduct handles POST /api/products
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	if !c.checkAdminRole(ctx) {
		return
	}

	var req requests.CreateProductRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	product, err := c.productService.Create(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// UpdateProduct handles PATCH /api/products/:id
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	if !c.checkAdminRole(ctx) {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req requests.UpdateProductRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	product, err := c.productService.Update(id, req)
	if err != nil {
		if errors.Is(err, utils.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// DeleteProduct handles DELETE /api/products/:id
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	if !c.checkAdminRole(ctx) {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = c.productService.Delete(id)
	if err != nil {
		if errors.Is(err, utils.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
