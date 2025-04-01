package models

import "github.com/prikshit/chameleon-privacy-module/internal/privacy"

type Server struct {
	PrivacyManager *privacy.PrivacyManager
}

type GenerateStealthAccountRequest struct {
	PubKeyHex string `json:"pub_key"`
}

type RecoverPrivKeyRequest struct {
	RecipientPrivKey string `json:"recipient_privkey"`
	EphemeralPubKey  string `json:"ephemeral_pubkey"`
}
