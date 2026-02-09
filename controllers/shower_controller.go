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

type ShowerController struct {
	showerService services.ShowerService
	userService   services.UserService
}

func NewShowerController(showerService services.ShowerService, userService services.UserService) *ShowerController {
	return &ShowerController{
		showerService: showerService,
		userService:   userService,
	}
}

// checkShowerOwnership validates if the authenticated user can access/modify the shower
func (c *ShowerController) checkShowerOwnership(ctx *gin.Context, showerID int) (bool, error) {
	userID, exists := ctx.Get("userID")
	if !exists {
		return false, errors.New("authentication required")
	}

	// Get user to check role
	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		return false, err
	}

	// Role 2+ (admin) has full access
	if user.Role >= 2 {
		return true, nil
	}

	// Role 0-1 must own the shower
	shower, err := c.showerService.GetByID(showerID)
	if err != nil {
		return false, err
	}

	if shower.HostID != uint(userID.(int)) {
		return false, utils.ErrUnauthorizedShowerAccess
	}

	return true, nil
}

// GetShower handles GET /api/showers/:id
func (c *ShowerController) GetShower(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shower ID"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	shower, err := c.showerService.GetByID(id)
	if err != nil {
		if errors.Is(err, utils.ErrShowerNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Shower not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get user to check permissions
	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is admin (Role 2+) or the host
	if user.Role < 2 && shower.HostID != uint(userID.(int)) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to access this shower"})
		return
	}

	ctx.JSON(http.StatusOK, shower)
}

// ListShowers handles GET /api/showers - Admin only with pagination
func (c *ShowerController) ListShowers(ctx *gin.Context) {
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

	// Only Role 2+ (admins) can list all showers
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to list showers"})
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

	showers, total, err := c.showerService.GetAllPaginated(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": showers,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetMyShowers handles GET /api/showers/me - Returns showers where user is the host
func (c *ShowerController) GetMyShowers(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	showers, err := c.showerService.GetByHostID(uint(userID.(int)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": showers})
}

// CreateShower handles POST /api/showers
func (c *ShowerController) CreateShower(ctx *gin.Context) {
	var req requests.CreateShowerRequest

	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Set HostID to the authenticated user's ID
	showerReq := requests.CreateShowerRequest{
		Guests:      req.Guests,
		ShowerDate:  req.ShowerDate,
		WeddingDate: req.WeddingDate,
		Location:    req.Location,
	}

	shower, err := c.showerService.CreateWithHost(showerReq, uint(userID.(int)))
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Host user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, shower)
}

// UpdateShower handles PATCH /api/showers/:id
func (c *ShowerController) UpdateShower(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shower ID"})
		return
	}

	// Check ownership
	authorized, err := c.checkShowerOwnership(ctx, id)
	if err != nil || !authorized {
		if errors.Is(err, utils.ErrUnauthorizedShowerAccess) || !authorized {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this shower"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req requests.UpdateShowerRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	shower, err := c.showerService.Update(id, req)
	if err != nil {
		if errors.Is(err, utils.ErrShowerNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Shower not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, shower)
}

// AddCatalog handles POST /api/showers/:id/catalog
func (c *ShowerController) AddCatalog(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shower ID"})
		return
	}

	// Check ownership
	authorized, err := c.checkShowerOwnership(ctx, id)
	if err != nil || !authorized {
		if errors.Is(err, utils.ErrUnauthorizedShowerAccess) || !authorized {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this shower"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req requests.AddCatalogRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	shower, err := c.showerService.AddCatalog(id, req)
	if err != nil {
		if errors.Is(err, utils.ErrShowerNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Shower not found"})
			return
		}
		if errors.Is(err, utils.ErrCatalogAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Catalog already exists for this shower"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, shower)
}

// AddPreferences handles POST /api/showers/:id/preferences
func (c *ShowerController) AddPreferences(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shower ID"})
		return
	}

	// Check ownership
	authorized, err := c.checkShowerOwnership(ctx, id)
	if err != nil || !authorized {
		if errors.Is(err, utils.ErrUnauthorizedShowerAccess) || !authorized {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this shower"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req requests.AddPreferencesRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	shower, err := c.showerService.AddPreferences(id, req)
	if err != nil {
		if errors.Is(err, utils.ErrShowerNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Shower not found"})
			return
		}
		if errors.Is(err, utils.ErrPreferencesAlreadyExist) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Preferences already exist for this shower"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, shower)
}

// GetShowerCatalog handles GET /api/showers/:id/catalog
func (c *ShowerController) GetShowerCatalog(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shower ID"})
		return
	}

	// Check ownership
	authorized, err := c.checkShowerOwnership(ctx, id)
	if err != nil || !authorized {
		if errors.Is(err, utils.ErrUnauthorizedShowerAccess) || !authorized {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to access this shower"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	catalog, products, err := c.showerService.GetCatalogWithProducts(id)
	if err != nil {
		if errors.Is(err, utils.ErrShowerNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Shower not found"})
			return
		}
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found for this shower"})
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

// GetAdminDashboard handles GET /api/admin/dashboard
func (c *ShowerController) GetAdminDashboard(ctx *gin.Context) {
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

	// Only Role 2+ (admins) can access dashboard
	if user.Role < 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to access admin dashboard"})
		return
	}

	stats, err := c.showerService.GetDashboardStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}
