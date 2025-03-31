package privacy

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/prikshit/chameleon-privacy-module/internal/sanctions"

	"github.com/ethereum/go-ethereum/crypto"
)

var ErrSanctionedAddress = errors.New("address is sanctioned")

// PrivacyManager manages stealth address generation and sanction detection.
type PrivacyManager struct {
	Detector *sanctions.Detector
}

// NewPrivacyManager creates a new PrivacyManager instance.
func NewPrivacyManager(detector *sanctions.Detector) *PrivacyManager {
	return &PrivacyManager{Detector: detector}
}

// GenerateStealthAddress generates a stealth address using the recipient's public key.
func (pm *PrivacyManager) GenerateStealthAddress(pubKey *ecdsa.PublicKey) (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
	// Check if the public key is sanctioned
	address := crypto.PubkeyToAddress(*pubKey).Hex()
	if pm.Detector.IsSanctioned(address) {
		return nil, nil, ErrSanctionedAddress
	}

	// Generate a random ephemeral key
	ephemeralPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Calculate the stealth address using ECDH (Elliptic Curve Diffie-Hellman)
	tempX, tempY := crypto.S256().Add(pubKey.X, pubKey.Y, ephemeralPrivKey.PublicKey.X, ephemeralPrivKey.PublicKey.Y)

	stealthPub := &ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     tempX,
		Y:     tempY,
	}

	// Return both the generated stealth address and the ephemeral private key for recovery.
	return stealthPub, ephemeralPrivKey, nil
}

// GenerateSharedSecret generates a shared secret using the recipient's private key and ephemeral public key.
func (pm *PrivacyManager) GenerateSharedSecret(privKey *ecdsa.PrivateKey, ephemeralPub *ecdsa.PublicKey) ([]byte, error) {
	// Perform ECDH: Shared secret = ephemeralPub * privKey
	x, _ := ephemeralPub.ScalarMult(ephemeralPub.X, ephemeralPub.Y, privKey.D.Bytes())
	secret := sha256.Sum256(x.Bytes())
	return secret[:], nil
}

// RecoverStealthPrivateKey recovers the recipient's stealth private key using their original private key and the ephemeral public key.
func (pm *PrivacyManager) RecoverStealthPrivateKey(recipientPriv *ecdsa.PrivateKey, ephemeralPub *ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	// Derive the shared secret
	sharedSecret, err := pm.GenerateSharedSecret(recipientPriv, ephemeralPub)
	if err != nil {
		return nil, err
	}

	// Convert shared secret to a big.Int
	sharedSecretInt := new(big.Int).SetBytes(sharedSecret)

	// Recover the stealth private key: stealthPriv = recipientPriv + sharedSecret
	stealthPriv := new(big.Int).Add(recipientPriv.D, sharedSecretInt)
	stealthPriv.Mod(stealthPriv, crypto.S256().Params().N)

	// Create a new ECDSA private key
	stealthPrivateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     ephemeralPub.X,
			Y:     ephemeralPub.Y,
		},
		D: stealthPriv,
	}

	return stealthPrivateKey, nil
}
