package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prikshit/chameleon-privacy-module/models"
)

// Verifies the Sanctioned Addresses
func HandleCheckSanction(c *gin.Context, s *models.Server) {
	var req struct {
		Address string `json:"address"`
	}

	// Bind JSON body to request struct using Gin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Sanction check based on the provided address
	isSanctioned := s.PrivacyManager.Detector.IsSanctioned(req.Address)

	// Return the result as a JSON response
	c.JSON(200, gin.H{"sanctioned": isSanctioned})
}
