package privacy

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
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

	// Convert Private Keys to Hex
	recipientPrivHex := fmt.Sprintf("0x%x", recipientPrivKey.D)
	fmt.Println("recipientPrivKey: ", recipientPrivHex)

	// Generate a stealth address using recipient's public key
	stealthPub, ephemeralPrivKey, err := pm.GenerateStealthAddress(&recipientPrivKey.PublicKey)
	assert.NoError(t, err)
	assert.NotNil(t, stealthPub)
	assert.NotNil(t, ephemeralPrivKey)
	// Convert Private Keys to Hex
	stealthPrivHex := fmt.Sprintf("0x%x", ephemeralPrivKey.D)
	fmt.Println("ephemeralPrivKey from Test: ", stealthPrivHex)

	// Recover the stealth private key using recipient's private key and the ephemeral public key
	recoveredPrivKey, err := pm.RecoverStealthPrivateKey(recipientPrivKey, &ephemeralPrivKey.PublicKey)
	assert.NoError(t, err)
	assert.NotNil(t, recoveredPrivKey)

	recoveredPrivHex := fmt.Sprintf("0x%x", recoveredPrivKey.D)
	fmt.Println("Recovered PrivKey from Test: ", recoveredPrivHex)

	// Validate recovered private key matches expected stealth key
	expectedStealthPub := &recoveredPrivKey.PublicKey
	assert.Equal(t, stealthPub.X, expectedStealthPub.X, "Recovered X-coord mismatch")
	assert.Equal(t, stealthPub.Y, expectedStealthPub.Y, "Recovered Y-coord mismatch")

	// Verify the algorithm used to recover the private Key is correct
	// Compute expected stealth private key using modular arithmetic
	curveN := recipientPrivKey.Curve.Params().N
	sharedSecretBytes, err := pm.GenerateSharedSecret(recipientPrivKey, &ephemeralPrivKey.PublicKey)
	assert.NoError(t, err)

	// Convert shared secret from []byte to *big.Int
	sharedSecret := new(big.Int).SetBytes(sharedSecretBytes)
	sharedSecret.Mod(sharedSecret, curveN) // Ensure it's within valid range

	expectedPrivKey := new(big.Int).Add(recipientPrivKey.D, sharedSecret)
	expectedPrivKey.Mod(expectedPrivKey, curveN) // Apply modulo n

	// Convert expected private key to Hex for debugging
	expectedPrivHex := fmt.Sprintf("0x%x", expectedPrivKey)
	fmt.Println("Expected Stealth Private Key (d_s):", expectedPrivHex)

	// Validate recovered private key matches expected stealth key modulo n
	assert.Equal(t, expectedPrivKey, recoveredPrivKey.D, "Recovered private key mismatch (mod n)")
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
