package network

import (
	"context"
	"fmt"
	"net"
	"regexp"

	"github.com/digitalocean/go-libvirt"
)

var macRe = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)

// libvirt-network.h (go-libvirt does not expose all enums)
const (
	networkUpdateSectionIPDhcpHost uint32 = 4 // VIR_NETWORK_UPDATE_SECTION_IP_DHCP_HOST
	networkUpdateCommandModify     uint32 = 2 // VIR_NETWORK_UPDATE_COMMAND_MODIFY
)

// SetStaticMapping ensures a DHCP reservation (MAC → IP) exists on a libvirt network.
//
// Semantics:
// - If mapping exists: update it to the requested IP (Modify; fallback Delete+Add)
// - If mapping does not exist: create it
func (m *LibvirtNetworkManager) SetStaticMapping(
	ctx context.Context,
	networkName, ip, mac string,
) error {
	if err := validateIP(ip); err != nil {
		return err
	}
	if err := validateMAC(mac); err != nil {
		return err
	}

	nw, err := m.conn.NetworkLookupByName(networkName)
	if err != nil {
		return fmt.Errorf("lookup network %q: %w", networkName, err)
	}

	flags := libvirt.NetworkUpdateAffectLive | libvirt.NetworkUpdateAffectConfig

	// 1) Modify (clean path if host entry exists)
	if err := m.modifyDHCPHost(nw, mac, ip, flags); err == nil {
		return nil
	}

	// 2) Fallback: Delete (ignore if missing) then Add
	// Delete by MAC selector only.
	_ = m.deleteDHCPHost(nw, mac, flags)

	if err := m.addDHCPHost(nw, mac, ip, flags); err != nil {
		return fmt.Errorf(
			"set dhcp mapping on network %q (mac=%s ip=%s): %w",
			networkName,
			mac,
			ip,
			err,
		)
	}

	return nil
}

func (m *LibvirtNetworkManager) modifyDHCPHost(
	nw libvirt.Network,
	mac, ip string,
	flags libvirt.NetworkUpdateFlags,
) error {
	xml := dhcpHostXML(mac, ip)

	return m.conn.NetworkUpdate(
		nw,
		networkUpdateCommandModify,
		networkUpdateSectionIPDhcpHost,
		-1,
		xml,
		flags,
	)
}

func (m *LibvirtNetworkManager) deleteDHCPHost(
	nw libvirt.Network,
	mac string,
	flags libvirt.NetworkUpdateFlags,
) error {
	selector := dhcpHostSelectorXML(mac)

	return m.conn.NetworkUpdate(
		nw,
		uint32(libvirt.NetworkUpdateCommandDelete),
		networkUpdateSectionIPDhcpHost,
		-1,
		selector,
		flags,
	)
}

func (m *LibvirtNetworkManager) addDHCPHost(
	nw libvirt.Network,
	mac, ip string,
	flags libvirt.NetworkUpdateFlags,
) error {
	xml := dhcpHostXML(mac, ip)

	return m.conn.NetworkUpdate(
		nw,
		uint32(libvirt.NetworkUpdateCommandAddLast),
		networkUpdateSectionIPDhcpHost,
		-1,
		xml,
		flags,
	)
}

func dhcpHostXML(mac, ip string) string {
	return fmt.Sprintf(`<host mac='%s' ip='%s'/>`, mac, ip)
}

func dhcpHostSelectorXML(mac string) string {
	return fmt.Sprintf(`<host mac='%s'/>`, mac)
}

func validateIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address %q", ip)
	}
	return nil
}

func validateMAC(mac string) error {
	if !macRe.MatchString(mac) {
		return fmt.Errorf("invalid MAC address %q", mac)
	}
	return nil
}
