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
	log.Println("Initializing server with PrivacyManager")
	return &models.Server{PrivacyManager: pm}
}

func Start(s *models.Server) error {
	r := gin.Default()

	// Register routes and pass the Server instance
	r.POST("/check-sanction", func(c *gin.Context) {
		log.Println("Handling check sanction request")
		controller.HandleCheckSanction(c, s)
	})

	r.POST("/generate-stealth", func(c *gin.Context) {
		log.Println("Handling generate stealth account request")
		controller.GenerateStealthAccount(c, s)
	})

	r.GET("/generate-account", func(c *gin.Context) {
		log.Println("Handling generate account request")
		controller.GenerateAccount(c)
	})
	r.POST("/recover-stealth-priv-key", func(c *gin.Context) {
		log.Println("Handling recover stealth private key request")
		controller.RecoverStealthPrivKey(c, s)
	})
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port: %s\n", port)
	return r.Run(":" + port)
}
