package network

import (
	"fmt"
	"net"
)

func MacFromIP(macPrefix, ipStr string) (string, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return "", fmt.Errorf("invalid IPv4 address %q", ipStr)
	}
	m, err := net.ParseMAC(macPrefix + ":00:00:00")
	if err != nil {
		return "", fmt.Errorf("invalid MAC prefix %q", macPrefix)
	}
	return fmt.Sprintf(
		"%02x:%02x:%02x:00:%02x:%02x",
		m[0], m[1], m[2],
		ip[2], ip[3],
	), nil
}

func ResolveMAC(prefix, ip, specMAC string) (string, error) {
	// Explicit MAC always wins
	if specMAC != "" {
		return specMAC, nil
	}

	// No MAC, no IP → nothing to resolve
	if ip == "" {
		return "", nil
	}

	mac, err := MacFromIP(prefix, ip)
	if err != nil {
		return "", err
	}

	return mac, nil
}
