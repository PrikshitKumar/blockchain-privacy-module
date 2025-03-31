package api

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prikshit/chameleon-privacy-module/internal/privacy"
)

type Server struct {
	PrivacyManager *privacy.PrivacyManager
}

func NewServer(pm *privacy.PrivacyManager) *Server {
	return &Server{PrivacyManager: pm}
}

func (s *Server) HandleCheckSanction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	isSanctioned := s.PrivacyManager.Detector.IsSanctioned(req.Address)
	json.NewEncoder(w).Encode(map[string]bool{"sanctioned": isSanctioned})
}

func (s *Server) HandleGenerateStealth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PubKeyHex string `json:"pub_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pubKey, err := crypto.UnmarshalPubkey(common.FromHex(req.PubKeyHex))
	if err != nil {
		http.Error(w, "Invalid public key", http.StatusBadRequest)
		return
	}

	stealthPub, stealthPriv, err := s.PrivacyManager.GenerateStealthAddress(pubKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	resp := map[string]string{
		"stealth_pub_key":  crypto.PubkeyToAddress(*stealthPub).Hex(),
		"stealth_priv_key": common.Bytes2Hex(crypto.FromECDSA(stealthPriv)),
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) Start() error {
	http.HandleFunc("/check-sanction", s.HandleCheckSanction)
	http.HandleFunc("/generate-stealth", s.HandleGenerateStealth)
	return http.ListenAndServe(":8080", nil)
}
