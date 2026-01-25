package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	RegisterUserRoutes(router)
	RegisterAuthRoutes(router)
	RegisterShowerRoutes(router)
	RegisterProductRoutes(router)
	RegisterCatalogRoutes(router)
	RegisterCatalogProductRoutes(router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}
