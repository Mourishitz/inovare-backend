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

type CatalogController struct {
	catalogService services.CatalogService
	userService    services.UserService
}

func NewCatalogController(catalogService services.CatalogService, userService services.UserService) *CatalogController {
	return &CatalogController{
		catalogService: catalogService,
		userService:    userService,
	}
}

// GetByURL handles GET /api/catalogs/url/:url
func (c *CatalogController) GetByURL(ctx *gin.Context) {
	url := ctx.Param("url")

	catalog, products, err := c.catalogService.GetProductsByURL(url)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		if errors.Is(err, utils.ErrCatalogNotApproved) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Catalog has not been approved yet"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"catalog":  catalog,
		"products": products,
	})
}

// PurchasePublicCatalog handles POST /api/public-catalogs/:slug/purchase.
func (c *CatalogController) PurchasePublicCatalog(ctx *gin.Context) {
	slug := ctx.Param("slug")

	var req requests.PublicCatalogPurchaseRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	response, err := c.catalogService.CreatePublicPurchase(ctx.Request.Context(), slug, req)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidProductID):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		case errors.Is(err, utils.ErrInvalidCustomerData):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Customer name, email, cellphone, and taxId are required"})
		case errors.Is(err, utils.ErrInvalidCurrentURL):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Current URL must be a valid absolute http(s) URL"})
		case errors.Is(err, utils.ErrInvalidTaxID):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "CPF ou CNPJ inválido"})
		case errors.Is(err, utils.ErrCatalogNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
		case errors.Is(err, utils.ErrCatalogProductNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		case errors.Is(err, utils.ErrCatalogProductAlreadyBought):
			ctx.JSON(http.StatusConflict, gin.H{"error": "Product is no longer available"})
		case errors.Is(err, utils.ErrCatalogNotApproved):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Catalog has not been approved yet"})
		case errors.Is(err, utils.ErrPaymentConfigurationMissing):
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Payment configuration is missing"})
		case errors.Is(err, utils.ErrPaymentProviderUnavailable):
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "Unable to create checkout right now"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetByID handles GET /api/catalogs/:id
func (c *CatalogController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog ID"})
		return
	}

	catalog, err := c.catalogService.GetByID(id)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, catalog)
}

// ApproveCatalog handles PATCH /api/catalogs/:id/approve
func (c *CatalogController) ApproveCatalog(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
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

	// Approve catalog
	catalog, err := c.catalogService.Approve(id, userID.(int))
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		if errors.Is(err, utils.ErrShowerNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Shower not found"})
			return
		}
		if errors.Is(err, utils.ErrUnauthorizedShowerAccess) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to approve this catalog"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, catalog)
}

// RegisterChanges handles PATCH /api/catalogs/:id/changes-made
func (c *CatalogController) RegisterChanges(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
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

	// Only admins (Role 2+) can register changes
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	catalog, err := c.catalogService.RegisterChanges(id)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, catalog)
}
