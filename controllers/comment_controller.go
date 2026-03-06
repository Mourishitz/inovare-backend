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

type CommentController struct {
	commentService services.CommentService
}

func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

// AddComment handles POST /api/catalogs/:id/comments
func (c *CommentController) AddComment(ctx *gin.Context) {
	catalogID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog ID"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req requests.CreateCommentRequest
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	comment, err := c.commentService.AddComment(catalogID, uint(userID.(int)), req.Content)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		if errors.Is(err, utils.ErrShowerNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Shower not found for this catalog"})
			return
		}
		if errors.Is(err, utils.ErrUnauthorizedShowerAccess) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Only the shower host can add comments"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

// ListComments handles GET /api/catalogs/:id/comments
func (c *CommentController) ListComments(ctx *gin.Context) {
	catalogID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid catalog ID"})
		return
	}

	comments, err := c.commentService.ListComments(catalogID)
	if err != nil {
		if errors.Is(err, utils.ErrCatalogNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Catalog not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
