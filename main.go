package main

import (
	"log"

	"inovare-backend/config"
	"inovare-backend/database"
	"inovare-backend/routes"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.GetConfig().ServerMode)

	database.Connect()
	router := gin.Default()

	routes.RegisterRoutes(router)

	port := config.GetConfig().ServerPort
	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
