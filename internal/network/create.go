package network

import (
	// "context"
	// "database/sql"
	"context"
	"fmt"

	// db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/vms"
)

// Create defines a Network in libvirt and inserts its database record.
// Create defines a Network in libvirt and inserts its database record.
func (m *LibvirtNetworkManager) Create(ctx context.Context, spec Config) error {
	// Ensure Libvirt connection is initialized
	if m.conn == nil {
		return vms.ErrNilLibvirtConn
	}

	// Validate that we have a network name
	name := spec.Name
	if name == "" {
		return ErrNetworkNameEmpty
	}

	// Prepare the database record
	// NOTE: NewNetworkRecord currently expects *Network. We might need to update it or create a temporary Network struct.
	// For now, let's look at how NewNetworkRecord works. It's likely in helpers.go or database package.
	// Assuming we can construct what it needs.
	// Let's defer NewNetworkRecord call update until we see helpers.go, but for now I'll create a dummy network struct if needed or update the helper.
	// Actually, I should update the helper to take Spec too.

	// Generate the network XML definition
	xmlConfig, err := m.prepareNetwork(spec)
	if err != nil {
		return fmt.Errorf("failed to prepare XML for network %q: %w", name, err)
	}

	// Define and start the network in libvirt
	if err := m.defineAndStartNetwork(xmlConfig); err != nil {
		return fmt.Errorf("failed to define/start network %q: %w", name, err)
	}

	// Insert into DB *after* successful libvirt creation (or before? original code did before).
	// Original code: record.Insert then define.
	// Let's stick to original order.

	// We need to refactor NewNetworkRecord. Let's assume we change it to accept spec.
	network := &Network{Spec: spec}
	record := NewNetworkRecord(network) // This needs to stick around or be updated?
	if err := record.Insert(ctx, m.db); err != nil {
		return fmt.Errorf("failed to insert database record for network %q: %w", name, err)
	}

	fmt.Printf("network/%s created\n", name)
	return nil
}
