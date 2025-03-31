package privacy

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
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

	// Generate ephemeral keypair
	ephemeralPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Compute shared secret: s = H(d_e * P_r)
	sharedX, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, ephemeralPrivKey.D.Bytes())
	sharedSecret := crypto.Keccak256(sharedX.Bytes()) // Hash for better randomness

	// Convert shared secret into scalar value
	s := new(big.Int).SetBytes(sharedSecret)
	fmt.Println("Shared Secret from Generation (s):", s.Text(16))

	// Compute stealth public key: P_s = P_r + s * G
	sGx, sGy := pubKey.Curve.ScalarBaseMult(s.Bytes())                         // s * G
	stealthPubX, stealthPubY := pubKey.Curve.Add(pubKey.X, pubKey.Y, sGx, sGy) // P_s = P_r + s * G

	stealthPub := &ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     stealthPubX,
		Y:     stealthPubY,
	}

	// Return stealth public key and ephemeral private key (needed for recipient to recover stealth private key)
	return stealthPub, ephemeralPrivKey, nil
}

// GenerateSharedSecret generates a shared secret using the recipient's private key and ephemeral public key.
func (pm *PrivacyManager) GenerateSharedSecret(privKey *ecdsa.PrivateKey, ephemeralPub *ecdsa.PublicKey) ([]byte, error) {
	// Perform ECDH: Shared secret = ephemeralPub * privKey
	sharedX, _ := privKey.Curve.ScalarMult(ephemeralPub.X, ephemeralPub.Y, privKey.D.Bytes())
	sharedSecret := crypto.Keccak256(sharedX.Bytes()) // Hash for randomness

	return sharedSecret, nil
}

// RecoverStealthPrivateKey recovers the recipient's stealth private key using their original private key and the ephemeral public key.
func (pm *PrivacyManager) RecoverStealthPrivateKey(recipientPriv *ecdsa.PrivateKey, ephemeralPub *ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	// Compute shared secret: s = H(d_r * P_e)
	sharedSecret, err := pm.GenerateSharedSecret(recipientPriv, ephemeralPub)
	if err != nil {
		return nil, err
	}

	// Convert shared secret into scalar value
	s := new(big.Int).SetBytes(sharedSecret)

	// Debug: Log values to verify correctness
	fmt.Println("Recipient Private Key (d_r):", recipientPriv.D.Text(16))
	fmt.Println("Shared Secret from Recovery (s):", s.Text(16))

	// Compute stealth private key: d_s = d_r + s mod n
	stealthPrivKey := new(big.Int).Add(recipientPriv.D, s)
	stealthPrivKey.Mod(stealthPrivKey, recipientPriv.Curve.Params().N) // Modulo to keep it in range

	// Debug: Verify stealth private key
	fmt.Println("Recovered Stealth Private Key (d_s):", stealthPrivKey.Text(16))

	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: recipientPriv.Curve,
			X:     recipientPriv.PublicKey.X,
			Y:     recipientPriv.PublicKey.Y,
		},
		D: stealthPrivKey,
	}, nil
}
