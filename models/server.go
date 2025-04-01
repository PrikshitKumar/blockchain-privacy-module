package models

import "github.com/prikshit/chameleon-privacy-module/internal/privacy"

type Server struct {
	PrivacyManager *privacy.PrivacyManager
}
