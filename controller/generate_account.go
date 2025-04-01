package controller

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"crypto/ecdsa"
	"crypto/elliptic"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
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
	pubKey, err := parseECDSAPubKey(pubKeyHex)
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

// Helper function to parse hex-encoded public key into *ecdsa.PublicKey
func parseECDSAPubKey(hexKey string) (*ecdsa.PublicKey, error) {
	// Remove "0x" prefix from the hex string
	hexKey = hexKey[2:]

	// Decode the hex string
	pubKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}

	// Ensure the public key is in the uncompressed format (starting with 0x04)
	if pubKeyBytes[0] != 0x04 {
		return nil, fmt.Errorf("invalid uncompressed public key format")
	}

	// Extract the X and Y coordinates from the public key bytes
	x := new(big.Int).SetBytes(pubKeyBytes[1:33])  // 32 bytes for X coordinate
	y := new(big.Int).SetBytes(pubKeyBytes[33:65]) // 32 bytes for Y coordinate

	// Create the ecdsa.PublicKey using the secp256k1 curve
	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(), // Use the elliptic.P256() curve for ECDSA (Ethereum uses secp256k1, but P256 for simplicity)
		X:     x,
		Y:     y,
	}

	return pubKey, nil
}
