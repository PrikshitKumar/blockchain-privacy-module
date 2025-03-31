package sanctions

import (
	"sync"
)

// Detector is responsible for detecting if a given address is sanctioned.
type Detector struct {
	SanctionedAddresses map[string]struct{} // We can replace the storage to Database for scalability.
	mu                  sync.RWMutex
}

// NewDetector creates a new Detector instance with an initial list of sanctioned addresses.
func NewDetector(initialAddresses []string) *Detector {
	detector := &Detector{
		SanctionedAddresses: make(map[string]struct{}),
	}

	for _, addr := range initialAddresses {
		detector.AddAddress(addr)
	}

	return detector
}

// AddAddress adds a new address to the sanctioned list.
func (d *Detector) AddAddress(address string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.SanctionedAddresses[address] = struct{}{}
}

// RemoveAddress removes an address from the sanctioned list.
func (d *Detector) RemoveAddress(address string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.SanctionedAddresses, address)
}

// IsSanctioned checks if a given address is sanctioned.
func (d *Detector) IsSanctioned(address string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.SanctionedAddresses[address]
	return exists
}
