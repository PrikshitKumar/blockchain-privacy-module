package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/prikshit/chameleon-privacy-module/models"
)

func GenerateStealthAccount(c *gin.Context, s *models.Server) {
	var req struct {
		PubKeyHex string `json:"pub_key"`
	}
	// Bind JSON body to request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	pubKeyHex := strings.ToLower(req.PubKeyHex)
	if pubKeyHex == "" {
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
	stealthPub, stealthPriv, err := s.PrivacyManager.GenerateStealthAddress(pubKey)
	if err != nil {
		log.Printf("Error generating stealth address: %v", err)
		return
	}

	// Return the generated keys as response
	c.JSON(http.StatusOK, gin.H{
		"stealth_pub_key":    crypto.PubkeyToAddress(*stealthPub).Hex(),
		"ephemeral_priv_key": common.Bytes2Hex(crypto.FromECDSA(stealthPriv)),
	})
}
