package controller

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/prikshit/chameleon-privacy-module/models"
)

// Generates the Stealth Account (by Payer)
func GenerateStealthAccount(c *gin.Context, s *models.Server) {
	var req models.GenerateStealthAccountRequest
	// Bind JSON body to request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.PubKeyHex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing pub_key"})
		return
	}

	// Unmarshal public key from hex string
	pubKey, err := crypto.UnmarshalPubkey(common.FromHex(req.PubKeyHex))
	if err != nil {
		log.Printf("Error unmarshalling public key: %v", err)
		return
	}

	// Generate stealth address and private key
	stealthPub, ephemeralPriv, err := s.PrivacyManager.GenerateStealthAddress(pubKey)
	if err != nil {
		log.Printf("Error generating stealth address: %v", err)
		return
	}

	// Return the generated keys as response
	c.JSON(http.StatusOK, gin.H{
		"stealth_pub_key":    "0x" + hex.EncodeToString(crypto.FromECDSAPub(stealthPub)),
		"ephemeral_priv_key": "0x" + fmt.Sprintf("%064s", common.Bytes2Hex(crypto.FromECDSA(ephemeralPriv))),
	})
}
