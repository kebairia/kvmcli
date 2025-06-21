package network

import (
	"encoding/xml"
	"fmt"
	"time"

	db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/utils"
)

// prepareNetwork generates the libvirt-compatible XML configuration
// for a virtual network, using optional parameters like DHCP and bridge.
func (net *VirtualNetwork) prepareNetwork() (string, error) {
	var opts []utils.NetworkOption

	// Append DHCP config if defined in the YAML
	if net.Config.Spec.DHCP != nil {
		start, startOk := net.Config.Spec.DHCP["start"]
		end, endOk := net.Config.Spec.DHCP["end"]
		if startOk && endOk {
			opts = append(opts, utils.WithDHCP(start, end))
		}
	}

	// Append bridge name if provided
	if net.Config.Spec.Bridge != "" {
		opts = append(opts, utils.WithBridge(net.Config.Spec.Bridge))
	}

	// Create the network definition with all options
	network := utils.NewNetwork(
		net.Config.Metadata.Name,
		net.Config.Spec.Mode,
		net.Config.Spec.Network.Address,
		net.Config.Spec.Network.Netmask,
		net.Config.Spec.Autostart,
		opts...,
	)

	xmlConfig, err := network.GenerateXML()
	if err != nil {
		return "", fmt.Errorf(
			"failed to generate XML for network %s: %v",
			net.Config.Metadata.Name,
			err,
		)
	}

	return xml.Header + string(xmlConfig), nil
}

// defineAndStartNetwork defines and starts the virtual network using libvirt.
func (net *VirtualNetwork) defineAndStartNetwork(xmlConfig string) error {
	// Define the network from the generated XML
	netInstance, err := net.conn.NetworkDefineXML(xmlConfig)
	if err != nil {
		return fmt.Errorf("failed to define network %s: %v", net.Config.Metadata.Name, err)
	}

	// Start (create) the defined network
	if err := net.conn.NetworkCreate(netInstance); err != nil {
		return fmt.Errorf("failed to start network %s: %w", net.Config.Metadata.Name, err)
	}
	// Set our network to be autostarted
	if net.Config.Spec.Autostart {
		net.conn.NetworkSetAutostart(netInstance, 1)
	}

	return nil
}

func NewVirtualNetworkRecord(net *VirtualNetwork) *db.VirtualNetworkRecord {
	// Create network record out of infos
	return &db.VirtualNetworkRecord{
		Name:       net.Config.Metadata.Name,
		Namespace:  net.Config.Metadata.Namespace,
		Labels:     net.Config.Metadata.Labels,
		MacAddress: net.Config.Spec.MacAddress,
		Bridge:     net.Config.Spec.Bridge,
		Mode:       net.Config.Spec.Mode,
		NetAddress: net.Config.Spec.Network.Address,
		Netmask:    net.Config.Spec.Network.Netmask,
		DHCP:       net.Config.Spec.DHCP,
		Autostart:  net.Config.Spec.Autostart,
		CreatedAt:  time.Now(),
	}
}
