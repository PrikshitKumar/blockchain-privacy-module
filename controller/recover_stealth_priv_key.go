package controller

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/prikshit/blockchain-privacy-module/models"
)

// Recovers the Private Key which can regenarate the Public Key and Address (Created by Payer) (By Receipient)
func RecoverStealthPrivKey(c *gin.Context, s *models.Server) {
	log.Println("Received request to recover stealth private key")

	var req models.RecoverPrivKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	log.Println("Parsing recipient private key")

	// Convert recipient private key from hex
	recipientPrivBytes, err := hex.DecodeString(req.RecipientPrivKey[2:])
	if err != nil {
		log.Println("Failed to decode recipient private key:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient private key"})
		return
	}
	recipientPrivKey, err := crypto.ToECDSA(recipientPrivBytes)
	if err != nil {
		log.Println("Failed to parse recipient private key:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse recipient private key"})
		return
	}

	// Convert ephemeral public key from hex
	ephemeralPubBytes, err := hex.DecodeString(req.EphemeralPubKey[2:])
	if err != nil {
		log.Println("Failed to decode ephemeral public key:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ephemeral public key"})
		return
	}
	ephemeralPubKey, err := crypto.UnmarshalPubkey(ephemeralPubBytes)
	if err != nil {
		log.Println("Failed to parse ephemeral public key:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse ephemeral public key"})
		return
	}

	log.Println("Recovering stealth private key")

	// Recover stealth private key
	recoveredPrivKey, err := s.PrivacyManager.RecoverStealthPrivateKey(recipientPrivKey, ephemeralPubKey)
	if err != nil {
		log.Println("Error recovering stealth private key:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert recovered private key to hex
	recoveredPrivHex := fmt.Sprintf("0x%x", recoveredPrivKey.D)
	log.Println("Successfully recovered stealth private key")

	stealthPub := &recoveredPrivKey.PublicKey
	log.Println("expectedStealthPub from Test123: ", stealthPub)

	// Convert the recovered public key to the same hex format for verification
	recoveredPubKeyHex := "0x" + hex.EncodeToString(crypto.FromECDSAPub(stealthPub))
	recoveredAddress := crypto.PubkeyToAddress(*stealthPub).Hex()

	c.JSON(http.StatusOK, gin.H{
		"recovered_priv_key":     recoveredPrivKey,
		"recovered_priv_key_hex": recoveredPrivHex,
		"stealth_pub":            stealthPub,
		"recovered_pub_key_hex":  recoveredPubKeyHex,
		"recovered_address":      recoveredAddress,
	})
}

// VerifyStealthKeys verifies if two stealth public keys match
func VerifyStealthKeys(c *gin.Context, s *models.Server) {
	log.Println("Received request to verify stealth keys")

	var req struct {
		GeneratedStealthPubKey string `json:"generated_stealth_pub_key"`
		RecoveredStealthPubKey string `json:"recovered_stealth_pub_key"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Compare the two public keys
	match := req.GeneratedStealthPubKey == req.RecoveredStealthPubKey

	log.Printf("Generated: %s", req.GeneratedStealthPubKey)
	log.Printf("Recovered: %s", req.RecoveredStealthPubKey)
	log.Printf("Keys match: %t", match)

	c.JSON(http.StatusOK, gin.H{
		"match":                     match,
		"generated_stealth_pub_key": req.GeneratedStealthPubKey,
		"recovered_stealth_pub_key": req.RecoveredStealthPubKey,
	})
}
