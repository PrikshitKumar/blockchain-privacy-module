package main

import (
	"log"

	"github.com/prikshit/blockchain-privacy-module/internal/privacy"
	"github.com/prikshit/blockchain-privacy-module/internal/sanctions"
	"github.com/prikshit/blockchain-privacy-module/server"
)

func main() {
	// Initialize list of sanctioned addresses
	initialAddresses := []string{"0xAbc123", "0xDef456"}
	log.Printf("Initializing sanctions detector with %d initial addresses\n", len(initialAddresses))

	// Create sanctions detector and PrivacyManager
	detector := sanctions.NewDetector(initialAddresses)
	log.Println("Sanctions detector initialized")

	privacyManager := privacy.NewPrivacyManager(detector)
	log.Println("Privacy manager initialized")

	// Initialize and start the server
	s := server.NewServer(privacyManager)
	log.Println("Server instance created")

	log.Println("Server starting on port 8080")
	if err := server.Start(s); err != nil {
		log.Fatal("Error starting the server: ", err)
	} else {
		log.Println("Server started successfully")
	}
}
