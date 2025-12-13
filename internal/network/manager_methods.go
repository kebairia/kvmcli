package network

import (
	"context"
	"fmt"
	"net"
	"regexp"

	"github.com/digitalocean/go-libvirt"
)

var macRe = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)

// These values match libvirt headers (libvirt-network.h).
// go-libvirt does not expose them all.
const (
	networkUpdateSectionIPDhcpHost uint32 = 4 // VIR_NETWORK_UPDATE_SECTION_IP_DHCP_HOST
	networkUpdateCommandModify     uint32 = 2 // VIR_NETWORK_UPDATE_COMMAND_MODIFY
)

// SetStaticMapping ensures a DHCP reservation (MAC → IP) exists on a libvirt network.
// Behavior:
// - If entry doesn't exist: add it
// - If entry exists: update it (Modify; fallback Delete+Add)
// - If already correct: no-op (best-effort; see note below)
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

	newEntry := dhcpHostXML(mac, ip)
	selector := dhcpHostSelectorXML(mac)

	// 1) Try Modify first: clean path when entry already exists.
	if err := m.conn.NetworkUpdate(
		nw,
		networkUpdateCommandModify,
		networkUpdateSectionIPDhcpHost,
		-1,
		newEntry,
		libvirt.NetworkUpdateFlags(flags),
	); err == nil {
		return nil
	}

	// 2) If Modify didn't work, try AddLast: handles "doesn't exist yet".
	if err := m.conn.NetworkUpdate(
		nw,
		uint32(libvirt.NetworkUpdateCommandAddLast),
		networkUpdateSectionIPDhcpHost,
		-1,
		newEntry,
		libvirt.NetworkUpdateFlags(flags),
	); err == nil {
		return nil
	}

	// 3) Fallback: Delete+Add (works when entry exists but Modify semantics differ).
	if err := m.conn.NetworkUpdate(
		nw,
		uint32(libvirt.NetworkUpdateCommandDelete),
		networkUpdateSectionIPDhcpHost,
		-1,
		selector,
		libvirt.NetworkUpdateFlags(flags),
	); err != nil {
		return fmt.Errorf(
			"update dhcp mapping (delete old) on network %q (mac=%s): %w",
			networkName,
			mac,
			err,
		)
	}

	if err := m.conn.NetworkUpdate(
		nw,
		uint32(libvirt.NetworkUpdateCommandAddLast),
		networkUpdateSectionIPDhcpHost,
		-1,
		newEntry,
		libvirt.NetworkUpdateFlags(flags),
	); err != nil {
		return fmt.Errorf(
			"update dhcp mapping (add new) on network %q (mac=%s ip=%s): %w",
			networkName,
			mac,
			ip,
			err,
		)
	}

	return nil
}

func dhcpHostXML(mac, ip string) string {
	return fmt.Sprintf(`<host mac='%s' ip='%s'/>`, mac, ip)
}

func dhcpHostSelectorXML(mac string) string {
	// Selector used for delete: match by MAC.
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
