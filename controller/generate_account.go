package controller

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// Generate a new Account or Key pair
func GenerateAccount(c *gin.Context) {
	log.Println("Received request to generate a new Ethereum account")

	// Generate a new private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		log.Println("Error generating private key: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating private key"})
		return
	}

	log.Println("Successfully generated private key")

	// Get the public key from the private key
	publicKeyECDSA := privateKey.PublicKey

	// Generate the Ethereum address from the public key
	address := crypto.PubkeyToAddress(publicKeyECDSA)
	log.Printf("Generated Ethereum address: %s", address.Hex())

	// Convert the private key to hexadecimal format
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	log.Printf("Private key (hex): %s", privateKeyHex)

	// Convert the public key (X, Y coordinates) to uncompressed format (0x04 prefix)
	pubKeyBytes := append([]byte{0x04}, publicKeyECDSA.X.Bytes()...)
	pubKeyBytes = append(pubKeyBytes, publicKeyECDSA.Y.Bytes()...)

	// Convert the uncompressed public key to hex format
	pubKeyHex := fmt.Sprintf("0x%x", pubKeyBytes)
	log.Printf("Public key (hex): %s", pubKeyHex)

	// Return the private key, public key (uncompressed), publick key in ecdsa format and Ethereum address as a JSON response
	log.Println("Returning generated account details")
	c.JSON(http.StatusOK, gin.H{
		"private_key":      "0x" + privateKeyHex, // Private key in hexadecimal
		"public_key":       pubKeyHex,            // Uncompressed public key (secp256k1)
		"public_key_ecdsa": publicKeyECDSA,       // Public Key ECDSA
		"address":          address.Hex(),        // Ethereum address (Hex format)
	})
}
