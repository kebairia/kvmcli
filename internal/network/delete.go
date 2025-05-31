package network

import (
	"fmt"
	"log"

	db "github.com/kebairia/kvmcli/internal/database"
)

// Delete removes a VirtualNetwork resource from libvirt and cleans up its database record.
func (vn *VirtualNetwork) Delete() error {
	if vn.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}

	name := vn.Metadata.Name

	// Lookup the network by its name.
	network, err := vn.Conn.NetworkLookupByName(name)
	if err != nil {
		return fmt.Errorf("failed to find network %q: %w", name, err)
	}

	// Destroy the network.
	if err := vn.Conn.NetworkDestroy(network); err != nil {
		return fmt.Errorf("failed to destroy network %q: %w", name, err)
	}

	// Undefine the network.
	if err := vn.Conn.NetworkUndefine(network); err != nil {
		return fmt.Errorf("failed to undefine network %q: %w", name, err)
	}

	// Delete the network record from the database.
	if err := db.Delete(db.Ctx, db.DB, name, db.NetworksTable); err != nil {
		log.Fatalf("failed to delete record for network %q: %v", name, err)
	}

	fmt.Printf("%s/%s deleted\n", "network", name)
	return nil
}
