package controller

import (
	"encoding/hex"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/prikshit/blockchain-privacy-module/models"
)

// Generates the Stealth Account (by Payer)
func GenerateStealthAccount(c *gin.Context, s *models.Server) {
	log.Println("Received request to generate a stealth account")

	var req models.GenerateStealthAccountRequest
	// Bind JSON body to request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.PubKeyHex == "" {
		log.Println("Missing pub_key in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing pub_key"})
		return
	}

	log.Println("Received public key:", req.PubKeyHex)

	// Unmarshal public key from hex string
	pubKey, err := crypto.UnmarshalPubkey(common.FromHex(req.PubKeyHex))
	if err != nil {
		log.Printf("Error unmarshalling public key: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid public key format"})
		return
	}

	log.Println("Successfully unmarshalled public key")

	// Generate stealth address and private key
	stealthPub, ephemeralPriv, err := s.PrivacyManager.GenerateStealthAddress(pubKey)
	if err != nil {
		log.Printf("Error generating stealth address: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate stealth address"})
		return
	}

	log.Println("Successfully generated stealth address")

	// Convert keys to hex format
	stealthPubHex := "0x" + hex.EncodeToString(crypto.FromECDSAPub(stealthPub))
	ephemeralPubHex := "0x" + hex.EncodeToString(crypto.FromECDSAPub(&ephemeralPriv.PublicKey))

	log.Println("stealthPub from Sender (Debug - 1): ", stealthPub)

	log.Println("Returning stealth account details")

	// Return the generated keys as response
	c.JSON(http.StatusOK, gin.H{
		"stealth_pub_key":   stealthPubHex,
		"ephemeral_pub_key": ephemeralPubHex,
	})
}
