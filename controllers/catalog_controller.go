package controllers

import (
	"errors"
	"net/http"
	"strconv"

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

	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Only admins (Role 2+) can approve catalogs
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Approve catalog
	catalog, err := c.catalogService.Approve(id)
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
