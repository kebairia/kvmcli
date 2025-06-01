package network

import (
	"fmt"

	db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/vms"
)

// Create defines a VirtualNetwork in libvirt and inserts its database record.
func (vnet *VirtualNetwork) Create() error {
	// Ensure Libvirt connection is initialized
	if vnet.Conn == nil {
		return vms.ErrNilLibvirtConn
	}

	// Validate that we have a network name
	name := vnet.Metadata.Name
	if name == "" {
		return VirtualNetworkNameEmpty
	}

	// Prepare the database record
	record := NewVirtualNetworkRecord(vnet)
	if err := record.Insert(db.Ctx, db.DB); err != nil {
		return fmt.Errorf("failed to insert database record for network %q: %w", name, err)
	}

	// Generate the network XML definition
	xmlConfig, err := vnet.prepareNetwork()
	if err != nil {
		return fmt.Errorf("failed to prepare XML for network %q: %w", name, err)
	}

	// Define and start the network in libvirt
	if err := vnet.defineAndStartNetwork(xmlConfig); err != nil {
		return fmt.Errorf("failed to define/start network %q: %w", name, err)
	}

	fmt.Printf("network/%s created\n", name)
	return nil
}
