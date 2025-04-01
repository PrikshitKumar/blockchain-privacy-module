package helpers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// ParseECDSAPubKey converts a hex-encoded public key to an *ecdsa.PublicKey.
func ParseECDSAPubKey(hexKey string) (*ecdsa.PublicKey, error) {
	log.Printf("Received hex public key for parsing in ECDSA: %s", hexKey)

	// Ensure the hex key is valid and has a "0x" prefix
	if len(hexKey) < 2 || hexKey[:2] != "0x" {
		err := fmt.Errorf("public key must start with '0x'")
		log.Println("Error:", err)
		return nil, err
	}
	hexKey = hexKey[2:]

	// Decode the hex string
	pubKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Printf("Failed to decode hex string: %v", err)
		return nil, fmt.Errorf("failed to decode hex string: %v", err)
	}

	log.Println("Successfully decoded hex string to bytes")

	// Ensure the public key is in uncompressed format (0x04 prefix)
	if len(pubKeyBytes) != 65 || pubKeyBytes[0] != 0x04 {
		err := fmt.Errorf("invalid public key format: expected uncompressed key (65 bytes)")
		log.Println("Error:", err)
		return nil, err
	}

	log.Println("Public key format verified as uncompressed (0x04 prefix)")

	// Extract X and Y coordinates (32 bytes each)
	x := new(big.Int).SetBytes(pubKeyBytes[1:33])
	y := new(big.Int).SetBytes(pubKeyBytes[33:])

	log.Printf("Extracted X coordinate: %s", x.Text(16))
	log.Printf("Extracted Y coordinate: %s", y.Text(16))

	// Create the ecdsa.PublicKey using secp256k1 (correct curve for Ethereum)
	pubKey := &ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     x,
		Y:     y,
	}
	log.Println("Successfully created ECDSA public key using secp256k1 curve")

	return pubKey, nil
}
