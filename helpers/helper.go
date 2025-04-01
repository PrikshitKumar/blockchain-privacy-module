package helpers

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

func ParseECDSAPubKey(hexKey string) (*ecdsa.PublicKey, error) {
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
