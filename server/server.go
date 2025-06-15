package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prikshit/blockchain-privacy-module/controller"
	"github.com/prikshit/blockchain-privacy-module/internal/privacy"
	"github.com/prikshit/blockchain-privacy-module/models"
)

func NewServer(pm *privacy.PrivacyManager) *models.Server {
	log.Println("Initializing server with PrivacyManager")
	return &models.Server{PrivacyManager: pm}
}

func Start(s *models.Server) error {
	r := gin.Default()

	// Register routes and pass the Server instance
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

	r.POST("/verify-stealth-keys", func(c *gin.Context) {
		log.Println("Handling verify stealth keys request")
		controller.VerifyStealthKeys(c, s)
	})

	r.POST("/sanctions/add", func(c *gin.Context) {
		log.Println("Handling add sanction request")
		controller.HandleAddSanctionedAddress(c, s)
	})

	r.POST("/sanctions/remove", func(c *gin.Context) {
		log.Println("Handling add sanction request")
		controller.HandleRemoveSanctionedAddress(c, s)
	})

	r.POST("/sanctions/check", func(c *gin.Context) {
		log.Println("Handling add sanction request")
		controller.HandleCheckSanction(c, s)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port: %s\n", port)
	return r.Run(":" + port)
}
