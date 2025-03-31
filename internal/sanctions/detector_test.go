package sanctions

import (
	"testing"
)

func TestDetector(t *testing.T) {
	initialList := []string{
		"0xAbc123",
		"0xDef456",
	}

	detector := NewDetector(initialList)

	// Test adding a new address
	detector.AddAddress("0xNew789")
	if !detector.IsSanctioned("0xNew789") {
		t.Error("Address should be marked as sanctioned")
	}

	// Test removing an address
	detector.RemoveAddress("0xNew789")
	if detector.IsSanctioned("0xNew789") {
		t.Error("Address should not be marked as sanctioned")
	}

	// Test checking initial addresses
	if !detector.IsSanctioned("0xAbc123") {
		t.Error("Address should be marked as sanctioned")
	}
	if detector.IsSanctioned("0xUnknown") {
		t.Error("Address should not be marked as sanctioned")
	}
}
