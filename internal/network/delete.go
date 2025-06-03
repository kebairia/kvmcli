package network

import (
	"errors"
	"fmt"
)

// Delete removes a VirtualNetwork from libvirt and deletes its database record.
func (vnet *VirtualNetwork) Delete() error {
	// Ensure we have a Libvirt connection
	if vnet.Conn == nil {
		return errors.New("libvirt connection is not initialized")
	}

	name := vnet.Metadata.Name

	// Lookup the network by name
	virNet, err := vnet.Conn.NetworkLookupByName(name)
	if err != nil {
		return fmt.Errorf("network %q not found: %w", name, err)
	}

	// Destroy the network (stop it if it’s running)
	if err := vnet.Conn.NetworkDestroy(virNet); err != nil {
		return fmt.Errorf("failed to destroy network %q: %w", name, err)
	}

	// Undefine the network (remove its definition from libvirt)
	if err := vnet.Conn.NetworkUndefine(virNet); err != nil {
		return fmt.Errorf("failed to undefine network %q: %w", name, err)
	}

	// Remove the record from the database
	record := NewVirtualNetworkRecord(vnet)
	if err := record.Delete(vnet.Context, vnet.DB); err != nil {
		return fmt.Errorf("failed to delete database record for network %q: %w", name, err)
	}

	fmt.Printf("network/%s deleted\n", name)
	return nil
}
