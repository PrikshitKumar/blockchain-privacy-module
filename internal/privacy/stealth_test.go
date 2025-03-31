package privacy

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prikshit/chameleon-privacy-module/internal/sanctions"
	"github.com/stretchr/testify/assert"
)

func TestGenerateStealthAddress(t *testing.T) {
	// Create a new detector instance with an empty sanctions list.
	detector := sanctions.NewDetector(nil)
	pm := NewPrivacyManager(detector)

	// Generate recipient's key pair
	recipientPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)

	// Generate a stealth address
	stealthPub, stealthPriv, err := pm.GenerateStealthAddress(&recipientPrivKey.PublicKey)

	assert.NoError(t, err)
	assert.NotNil(t, stealthPub)
	assert.NotNil(t, stealthPriv)

	// Convert Stealth Public Key to Ethereum Address
	stealthAddress := crypto.PubkeyToAddress(*stealthPub).Hex()

	// Convert Private Keys to Hex
	recipientPrivHex := fmt.Sprintf("0x%x", recipientPrivKey.D)
	stealthPrivHex := fmt.Sprintf("0x%x", stealthPriv.D)

	// Convert Public Keys to Ethereum Addresses
	recipientPubAddress := crypto.PubkeyToAddress(recipientPrivKey.PublicKey).Hex()
	stealthPubAddress := crypto.PubkeyToAddress(stealthPriv.PublicKey).Hex()

	// Print the formatted keys
	fmt.Println("Recipient Private Key:", recipientPrivHex)
	fmt.Println("Recipient Public Address:", recipientPubAddress)
	fmt.Println("Stealth Private Key:", stealthPrivHex)
	fmt.Println("Stealth Public Address:", stealthPubAddress)
	fmt.Println("Stealth Address:", stealthAddress)
}

func TestGenerateSharedSecret(t *testing.T) {
	detector := sanctions.NewDetector(nil)
	pm := NewPrivacyManager(detector)

	// Generate recipient's key pair
	recipientPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)

	// Generate stealth key pair (simulating sender)
	stealthPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)

	// Generate shared secret
	sharedSecret, err := pm.GenerateSharedSecret(recipientPrivKey, &stealthPrivKey.PublicKey)
	assert.NoError(t, err)
	assert.NotNil(t, sharedSecret)
	assert.Len(t, sharedSecret, 32) // Expected SHA256 hash length

	// Convert shared secret to hex
	sharedSecretHex := fmt.Sprintf("%x", sharedSecret)

	// Print formatted output
	fmt.Println("Shared Secret:", sharedSecretHex)
}

func TestRecoverStealthPrivateKey(t *testing.T) {
	detector := sanctions.NewDetector(nil)
	pm := NewPrivacyManager(detector)

	// Generate recipient's key pair
	recipientPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)

	// Generate stealth key pair (simulating sender)
	stealthPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)

	// Generate a stealth address
	stealthPub, _, err := pm.GenerateStealthAddress(&recipientPrivKey.PublicKey)
	assert.NoError(t, err)
	assert.NotNil(t, stealthPub)

	// Recover the stealth private key
	recoveredPrivKey, err := pm.RecoverStealthPrivateKey(recipientPrivKey, &stealthPrivKey.PublicKey)
	assert.NoError(t, err)
	assert.NotNil(t, recoveredPrivKey)
	assert.Equal(t, recipientPrivKey.D, recoveredPrivKey.D, "The recovered private key should match the original recipient private key.")
}

func TestSanctionedAddress(t *testing.T) {
	// Create a detector with one sanctioned address
	detector := sanctions.NewDetector([]string{"0x1234567890abcdef1234567890abcdef12345678"})
	pm := NewPrivacyManager(detector)

	// Generate recipient's key pair
	recipientPrivKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	assert.NoError(t, err)

	// Mark the recipient address as sanctioned
	address := crypto.PubkeyToAddress(recipientPrivKey.PublicKey).Hex()
	detector.SanctionedAddresses[address] = struct{}{}

	// Attempt to generate a stealth address (should fail)
	_, _, err = pm.GenerateStealthAddress(&recipientPrivKey.PublicKey)
	assert.ErrorIs(t, err, ErrSanctionedAddress)
}
