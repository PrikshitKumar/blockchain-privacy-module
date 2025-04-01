package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prikshit/chameleon-privacy-module/controller"
	"github.com/prikshit/chameleon-privacy-module/internal/privacy"
	"github.com/prikshit/chameleon-privacy-module/models"
)

func NewServer(pm *privacy.PrivacyManager) *models.Server {
	return &models.Server{PrivacyManager: pm}
}

func Start(s *models.Server) error {
	r := gin.Default()

	// Register routes and pass the Server instance
	r.POST("/check-sanction", func(c *gin.Context) {
		controller.HandleCheckSanction(c, s)
	})

	r.POST("/generate-stealth", func(c *gin.Context) {
		controller.GenerateStealthAccount(c, s)
	})
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port:", port)
	return r.Run(":" + port)
}
