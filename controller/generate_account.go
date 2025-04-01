package controller

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// Generate a new Account or Key pair
func GenerateAccount(c *gin.Context) {
	// Generate a new private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
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

	// Return the private key, public key (uncompressed), publick key in ecdsa format and Ethereum address as a JSON response
	c.JSON(200, gin.H{
		"private_key":      "0x" + privateKeyHex, // Private key in hexadecimal
		"public_key":       pubKeyHex,            // Uncompressed public key (secp256k1)
		"public_key_ecdsa": publicKeyECDSA,       // Public Key ECDSA
		"address":          address.Hex(),        // Ethereum address (Hex format)
	})
}
