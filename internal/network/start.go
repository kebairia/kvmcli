package network

import (
	"context"
)

// Start starts the network.
func (m *LibvirtNetworkManager) Start(ctx context.Context, name string) error {
	// TODO: Implement start logic if separate from creation.
	// For now, Create also Starts.
	// If Start is called separately, we should lookup and create.
	return nil
}
