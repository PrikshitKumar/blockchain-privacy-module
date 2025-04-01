package controller

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/prikshit/chameleon-privacy-module/helpers"
)

func GenerateAccount(c *gin.Context) {
	// Generate a new private key using secp256k1
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Error generating private key: %v", err)
		c.JSON(500, gin.H{"error": "Error generating private key"})
		return
	}

	// Get the public key from the private key
	publicKeyECDSA := privateKey.PublicKey

	// Generate the Ethereum address from the public key
	address := crypto.PubkeyToAddress(publicKeyECDSA)

	// Convert the private key to hexadecimal format
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	// Convert the public key (X, Y coordinates) to uncompressed format (0x04 prefix)
	pubKeyBytes := append([]byte{0x04}, publicKeyECDSA.X.Bytes()...)
	pubKeyBytes = append(pubKeyBytes, publicKeyECDSA.Y.Bytes()...)

	// Convert the uncompressed public key to hex format
	pubKeyHex := fmt.Sprintf("0x%x", pubKeyBytes)

	// Convert the public key bytes back to ecdsa.PublicKey
	pubKey, err := helpers.ParseECDSAPubKey(pubKeyHex)
	if err != nil {
		log.Printf("Error parsing public key: %v", err)
		c.JSON(500, gin.H{"error": "Invalid public key format"})
		return
	}

	// Now, you can pass the pubKey to GenerateStealthAddress or any other function expecting an *ecdsa.PublicKey
	// Example: stealthAddress := pm.GenerateStealthAddress(pubKey)

	// Return the private key, public key (uncompressed), and Ethereum address as a JSON response
	c.JSON(200, gin.H{
		"private_key":       privateKeyHex, // Private key in hexadecimal
		"public_key":        pubKeyHex,     // Uncompressed public key (secp256k1)
		"public_key_ecdsa":  publicKeyECDSA,
		"address":           address.Hex(), // Ethereum address (Hex format)
		"parsed_public_key": pubKey,
	})
}
