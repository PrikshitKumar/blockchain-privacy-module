package controller

import (
	"encoding/hex"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/prikshit/chameleon-privacy-module/models"
)

func RecoverStealthPrivKey(c *gin.Context, s *models.Server) {
	var req models.RecoverPrivKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Convert recipient private key from hex
	recipientPrivBytes, err := hex.DecodeString(req.RecipientPrivKey[2:])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient private key"})
		return
	}
	recipientPrivKey, err := crypto.ToECDSA(recipientPrivBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse recipient private key"})
		return
	}

	// Convert ephemeral public key from hex
	ephemeralPubBytes, err := hex.DecodeString(req.EphemeralPubKey[2:])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ephemeral public key"})
		return
	}
	ephemeralPubKey, err := crypto.UnmarshalPubkey(ephemeralPubBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse ephemeral public key"})
		return
	}

	// Recover stealth private key
	recoveredPrivKey, err := s.PrivacyManager.RecoverStealthPrivateKey(recipientPrivKey, ephemeralPubKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert recovered private key to hex
	// recoveredPrivHex := fmt.Sprintf("0x%x", recoveredPrivKey.D)

	c.JSON(http.StatusOK, gin.H{"recovered_privkey": recoveredPrivKey})
}
