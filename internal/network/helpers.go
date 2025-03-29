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
	if net.Spec.DHCP != nil {
		start, startOk := net.Spec.DHCP["start"]
		end, endOk := net.Spec.DHCP["end"]
		if startOk && endOk {
			opts = append(opts, utils.WithDHCP(start, end))
		}
	}

	// Append bridge name if provided
	if net.Spec.Bridge != "" {
		opts = append(opts, utils.WithBridge(net.Spec.Bridge))
	}

	// Create the network definition with all options
	network := utils.NewNetwork(
		net.Metadata.Name,
		net.Spec.Mode,
		net.Spec.NetAddress,
		net.Spec.Netmask,
		net.Spec.Autostart,
		opts...,
	)

	xmlConfig, err := network.GenerateXML()
	if err != nil {
		return "", fmt.Errorf("failed to generate XML for network %s: %v", net.Metadata.Name, err)
	}

	return xml.Header + string(xmlConfig), nil
}

// defineAndStartNetwork defines and starts the virtual network using libvirt.
func (net *VirtualNetwork) defineAndStartNetwork(xmlConfig string) error {
	// Define the network from the generated XML
	netInstance, err := net.Conn.NetworkDefineXML(xmlConfig)
	if err != nil {
		return fmt.Errorf("failed to define network %s: %v", net.Metadata.Name, err)
	}

	// Start (create) the defined network
	if err := net.Conn.NetworkCreate(netInstance); err != nil {
		return fmt.Errorf("failed to start network %s: %w", net.Metadata.Name, err)
	}
	// Set our network to be autostarted
	if net.Spec.Autostart {
		net.Conn.NetworkSetAutostart(netInstance, 1)
	}

	return nil
}

func NewNetRecord(net *VirtualNetwork) *db.NetRecord {
	// Create vm record out of infos
	return &db.NetRecord{
		Name:       net.Metadata.Name,
		Namespace:  net.Metadata.Namespace,
		Labels:     net.Metadata.Labels,
		MacAddress: net.Spec.MacAddress,
		Bridge:     net.Spec.Bridge,
		Mode:       net.Spec.Mode,
		NetAddress: net.Spec.NetAddress,
		Netmask:    net.Spec.Netmask,
		DHCP:       net.Spec.DHCP,
		CreatedAt:  time.Now(),
	}
}
