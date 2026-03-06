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

type CatalogProductController struct {
	catalogProductService services.CatalogProductService
	userService           services.UserService
}

func NewCatalogProductController(
	catalogProductService services.CatalogProductService,
	userService services.UserService,
) *CatalogProductController {
	return &CatalogProductController{
		catalogProductService: catalogProductService,
		userService:           userService,
	}
}

// AttachProduct handles POST /api/catalogs/:id/attach-product
func (c *CatalogProductController) AttachProduct(ctx *gin.Context) {
	catalogID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog ID"})
		return
	}

	// Get authenticated user
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Only admins (Role 2+) can attach products to catalogs
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Parse request body
	var req requests.AttachProductToCatalogRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	// Attach product to catalog
	catalogProduct, err := c.catalogProductService.AttachProduct(catalogID, req)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		if errors.Is(err, utils.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		if errors.Is(err, utils.ErrProductAlreadyInCatalog) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Product already exists in this catalog"})
			return
		}
		if errors.Is(err, utils.ErrProductIsExclusive) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Product is exclusive and already assigned to another catalog"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, catalogProduct)
}

// ListCatalogProducts handles GET /api/catalogs/:id/products
func (c *CatalogProductController) ListCatalogProducts(ctx *gin.Context) {
	catalogID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog ID"})
		return
	}

	// Get authenticated user
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Only admins (Role 2+) can list catalog products
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// List catalog products
	catalogProducts, err := c.catalogProductService.ListCatalogProducts(catalogID)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, catalogProducts)
}

// UpdateCatalogProduct handles PATCH /api/catalogs/:id/update-product/:product_id
func (c *CatalogProductController) UpdateCatalogProduct(ctx *gin.Context) {
	_, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog ID"})
		return
	}

	productID, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog product ID"})
		return
	}

	// Get authenticated user
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Only admins (Role 2+) can update catalog products
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Parse request body
	var req requests.UpdateCatalogProductRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	// Update catalog product
	catalogProduct, err := c.catalogProductService.UpdateCatalogProduct(productID, req)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, catalogProduct)
}

// CreateExclusiveProduct handles POST /api/catalogs/:id/exclusive-products
func (c *CatalogProductController) CreateExclusiveProduct(ctx *gin.Context) {
	catalogID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog ID"})
		return
	}

	// Get authenticated user
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Only admins (Role 2+) can create exclusive products
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req requests.CreateExclusiveProductRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	catalogProduct, err := c.catalogProductService.CreateExclusiveProduct(catalogID, req)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, catalogProduct)
}

// DetachProduct handles DELETE /api/catalogs/:id/detach-product/:product_id
func (c *CatalogProductController) DetachProduct(ctx *gin.Context) {
	catalogID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog ID"})
		return
	}

	productID, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Get authenticated user
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Only admins (Role 2+) can detach products from catalogs
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Detach product from catalog
	err = c.catalogProductService.DetachProduct(catalogID, productID)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product detached from catalog successfully"})
}
