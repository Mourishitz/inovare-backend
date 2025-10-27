package controllers

import (
	"errors"
	"net/http"

	"inovare-backend/requests"
	"inovare-backend/services"
	"inovare-backend/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService services.UserService
}

func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req requests.LoginRequest

	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	user, err := c.userService.GetByEmail(req.Email)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrUserNotFound):
			ctx.JSON(http.StatusConflict, gin.H{"error": "user not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	if !utils.CheckValidHashWithPassword(req.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

func (c *AuthController) GetSessionUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := c.userService.GetByID(userID.(int))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
