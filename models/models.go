package models

import "github.com/prikshit/blockchain-privacy-module/internal/privacy"

type Server struct {
	PrivacyManager *privacy.PrivacyManager
}

type GenerateStealthAccountRequest struct {
	PubKeyHex string `json:"pub_key"`
}

type GenerateStealthAccountResponse struct {
	StealthPubKey   string `json:"stealth_pub_key"`
	EphemeralPubKey string `json:"ephemeral_pub_key"`
}

type RecoverPrivKeyRequest struct {
	RecipientPrivKey string `json:"recipient_privkey"`
	EphemeralPubKey  string `json:"ephemeral_pubkey"`
}
