package network

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
)

func (net *VirtualNetwork) Delete() error {
	if net.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}
	netName := net.Metadata.Name
	network, err := net.Conn.NetworkLookupByName(netName)
	if err != nil {
		return fmt.Errorf("Failed to find Network %s: %w", netName, err)
	}
	// Attempt to destroy the network.
	if err := net.Conn.NetworkDestroy(network); err != nil {
		return fmt.Errorf(
			"failed to detroy network %q: %w",
			netName,
			err,
		)
	}
	// Undefine the network
	if err := net.Conn.NetworkUndefine(network); err != nil {
		return fmt.Errorf("failed to undefine network %q: %w", netName, err)
	}
	err = database.DeleteNetwork(net.Metadata.Name)
	if err != nil {
		logger.Log.Errorf("failed to delete record for network %s: %v", net.Metadata.Name, err)
	}

	logger.Log.Infof("%s/%s deleted", "net", netName)
	return nil
}
