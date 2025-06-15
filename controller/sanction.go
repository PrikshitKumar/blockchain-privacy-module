package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prikshit/blockchain-privacy-module/models"
)

// HandleAddSanctionedAddress adds a new address to the sanctioned list.
func HandleAddSanctionedAddress(c *gin.Context, s *models.Server) {
	var request struct {
		Address string `json:"address" binding:"required"`
	}
	log.Println("Received request to add sanctioned address")

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("Adding address to sanctioned list: %s\n", request.Address)

	s.PrivacyManager.Detector.AddAddress(request.Address)

	log.Printf("Address %s added to sanctioned list successfully\n", request.Address)

	c.JSON(http.StatusOK, gin.H{"message": "Address added to sanctioned list"})
}

// HandleRemoveSanctionedAddress removes an address from the sanctioned list.
func HandleRemoveSanctionedAddress(c *gin.Context, s *models.Server) {
	var request struct {
		Address string `json:"address" binding:"required"`
	}

	log.Println("Received request to remove sanctioned address")

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("Removing address from sanctioned list: %s\n", request.Address)

	s.PrivacyManager.Detector.RemoveAddress(request.Address)

	log.Printf("Address %s removed from sanctioned list successfully\n", request.Address)

	c.JSON(http.StatusOK, gin.H{"message": "Address removed from sanctioned list"})
}

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
