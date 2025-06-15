package privacy

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"log"
	"math/big"

	"github.com/prikshit/blockchain-privacy-module/internal/sanctions"

	"github.com/ethereum/go-ethereum/crypto"
)

var ErrSanctionedAddress = errors.New("address is sanctioned")

// PrivacyManager manages stealth address generation and sanction detection.
type PrivacyManager struct {
	Detector *sanctions.Detector
}

// NewPrivacyManager creates a new PrivacyManager instance.
func NewPrivacyManager(detector *sanctions.Detector) *PrivacyManager {
	log.Println("Initializing PrivacyManager")
	return &PrivacyManager{Detector: detector}
}

// GenerateStealthAddress generates a stealth address using the recipient's public key.
func (pm *PrivacyManager) GenerateStealthAddress(pubKey *ecdsa.PublicKey) (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
	// Check if the public key is sanctioned
	address := crypto.PubkeyToAddress(*pubKey).Hex()
	log.Printf("Attempting to generate stealth address for: %s\n", address)

	if pm.Detector.IsSanctioned(address) {
		log.Printf("Sanctioned address detected: %s\n", address)
		return nil, nil, ErrSanctionedAddress
	}

	// Generate ephemeral keypair
	log.Println("Generating ephemeral keypair")
	ephemeralPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		log.Printf("Error generating ephemeral key: %v\n", err)
		return nil, nil, err
	}
	log.Println("Ephemeral keypair generated successfully")

	// Compute shared secret: s = H(d_e * P_r)
	log.Println("Computing shared secret")
	sharedX, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, ephemeralPrivKey.D.Bytes())
	sharedSecret := crypto.Keccak256(sharedX.Bytes()) // Hash for better randomness

	// Convert shared secret into scalar value
	s := new(big.Int).SetBytes(sharedSecret)
	log.Printf("Secret from Sender: %s\n", s)

	// Compute stealth public key: P_s = P_r + s * G
	sGx, sGy := pubKey.Curve.ScalarBaseMult(s.Bytes())                         // s * G
	stealthPubX, stealthPubY := pubKey.Curve.Add(pubKey.X, pubKey.Y, sGx, sGy) // P_s = P_r + s * G

	stealthPub := &ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     stealthPubX,
		Y:     stealthPubY,
	}

	log.Println("Stealth public key generated successfully")

	// Return stealth public key and ephemeral private key (needed for recipient to recover stealth private key)
	return stealthPub, ephemeralPrivKey, nil
}

// GenerateSharedSecret generates a shared secret using the recipient's private key and ephemeral public key.
func (pm *PrivacyManager) GenerateSharedSecret(privKey *ecdsa.PrivateKey, ephemeralPub *ecdsa.PublicKey) ([]byte, error) {
	log.Println("Generating shared secret using ECDH")

	// Perform ECDH: Shared secret = ephemeralPub * privKey
	sharedX, _ := privKey.Curve.ScalarMult(ephemeralPub.X, ephemeralPub.Y, privKey.D.Bytes())
	sharedSecret := crypto.Keccak256(sharedX.Bytes()) // Hash for randomness

	log.Printf("Shared secret generated: %x\n", sharedSecret)

	return sharedSecret, nil
}

// RecoverStealthPrivateKey recovers the recipient's stealth private key using their original private key and the ephemeral public key.
func (pm *PrivacyManager) RecoverStealthPrivateKey(recipientPriv *ecdsa.PrivateKey, ephemeralPub *ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	log.Println("Recovering stealth private key")

	// Compute shared secret: s = H(d_r * P_e)
	sharedSecret, err := pm.GenerateSharedSecret(recipientPriv, ephemeralPub)
	if err != nil {
		log.Printf("Error generating shared secret: %v\n", err)
		return nil, err
	}

	// Convert shared secret into scalar value
	s := new(big.Int).SetBytes(sharedSecret)
	log.Printf("Secret from Receiver during recovery: %s\n", s)

	// Compute stealth private key: d_s = (d_r + s) mod n
	log.Println("Computing stealth private key")
	stealthPrivKey := new(big.Int).Add(recipientPriv.D, s)
	stealthPrivKey.Mod(stealthPrivKey, recipientPriv.Curve.Params().N) // Modulo n to stay within valid range

	// Recompute public key from stealth private key
	stealthPubX, stealthPubY := recipientPriv.Curve.ScalarBaseMult(stealthPrivKey.Bytes()) // d_s * G
	stealthPublicKey := &ecdsa.PublicKey{
		Curve: recipientPriv.Curve,
		X:     stealthPubX,
		Y:     stealthPubY,
	}

	log.Println("Stealth private key recovered successfully")

	// Return stealth private key with corrected public key
	return &ecdsa.PrivateKey{
		PublicKey: *stealthPublicKey,
		D:         stealthPrivKey,
	}, nil
}
