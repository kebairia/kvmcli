package network

import (
	"encoding/xml"
	"fmt"
	"time"

	db "github.com/kebairia/kvmcli/internal/database"
	templates "github.com/kebairia/kvmcli/internal/templates"
)

// prepareNetwork generates the libvirt-compatible XML configuration
// for a virtual network, using optional parameters like DHCP and bridge.
func (net *Network) prepareNetwork() (string, error) {
	var opts []templates.NetworkOption

	// Append DHCP config if defined in the YAML
	// Append DHCP config if defined in the YAML
	if net.Spec.DHCP != nil {
		opts = append(opts, templates.WithDHCP(net.Spec.DHCP.Start, net.Spec.DHCP.End))
	}

	// Append bridge name if provided
	if net.Spec.Bridge != "" {
		opts = append(opts, templates.WithBridge(net.Spec.Bridge))
	}

	// Create the network definition with all options
	netXML := templates.NewNetwork(
		net.Spec.Name,
		net.Spec.Mode,
		net.Spec.NetAddress,
		net.Spec.NetMask,
		net.Spec.Autostart,
		opts...,
	)

	xmlConfig, err := netXML.GenerateXML()
	if err != nil {
		return "", fmt.Errorf(
			"failed to generate XML for network %s: %v",
			net.Spec.Name,
			err,
		)
	}

	return xml.Header + string(xmlConfig), nil
}

// defineAndStartNetwork defines and starts the virtual network using libvirt.
func (net *Network) defineAndStartNetwork(xmlConfig string) error {
	// Define the network from the generated XML
	netInstance, err := net.conn.NetworkDefineXML(xmlConfig)
	if err != nil {
		return fmt.Errorf("failed to define network %s: %v", net.Spec.Name, err)
	}

	// Start (create) the defined network
	if err := net.conn.NetworkCreate(netInstance); err != nil {
		return fmt.Errorf("failed to start network %s: %w", net.Spec.Name, err)
	}
	// Set our network to be autostarted
	if net.Spec.Autostart {
		net.conn.NetworkSetAutostart(netInstance, 1)
	}

	return nil
}

func NewNetworkRecord(net *Network) *db.VirtualNetwork {
	// Create a map[string]string for the DB record if DHCP is present
	var dhcpMap map[string]string
	if net.Spec.DHCP != nil {
		dhcpMap = map[string]string{
			"start": net.Spec.DHCP.Start,
			"end":   net.Spec.DHCP.End,
		}
	}

	// Create network record out of infos
	return &db.VirtualNetwork{
		Name:      net.Spec.Name,
		Namespace: net.Spec.Namespace,
		Labels:    net.Spec.Labels,
		// MacAddress: net.Spec.MacAddress,
		Bridge:     net.Spec.Bridge,
		Mode:       net.Spec.Mode,
		NetAddress: net.Spec.NetAddress,
		Netmask:    net.Spec.NetMask,
		DHCP:       dhcpMap,
		Autostart:  net.Spec.Autostart,
		CreatedAt:  time.Now(),
	}
}
