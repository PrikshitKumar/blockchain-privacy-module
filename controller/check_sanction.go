package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prikshit/chameleon-privacy-module/models"
)

// Verifies the Sanctioned Addresses
func HandleCheckSanction(c *gin.Context, s *models.Server) {
	var req struct {
		Address string `json:"address"`
	}

	log.Println("Received request to check sanction status")

	// Bind JSON body to request struct using Gin
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Checking sanction status for address: %s", req.Address)

	// Sanction check based on the provided address
	isSanctioned := s.PrivacyManager.Detector.IsSanctioned(req.Address)

	log.Printf("Sanction status for %s: %v", req.Address, isSanctioned)

	// Return the result as a JSON response
	c.JSON(http.StatusOK, gin.H{"sanctioned": isSanctioned})
}
