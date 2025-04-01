package main

import (
	"log"

	"github.com/prikshit/chameleon-privacy-module/internal/privacy"
	"github.com/prikshit/chameleon-privacy-module/internal/sanctions"
	"github.com/prikshit/chameleon-privacy-module/server"
)

func main() {
	initialAddresses := []string{"0xAbc123...", "0xDef456..."}
	detector := sanctions.NewDetector(initialAddresses)
	privacyManager := privacy.NewPrivacyManager(detector)
	s := server.NewServer(privacyManager)

	log.Println("Server starting on port 8080...")
	if err := server.Start(s); err != nil {
		log.Fatal(err)
	}
}
