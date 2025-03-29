package utils

import (
	"encoding/xml"
)

// Network represents the overall XML structure
// The Bridge/Forward field is a pointer, so if it's nil, it won't be included in the XML
type Network struct {
	XMLName xml.Name `xml:"network"`
	Name    string   `xml:"name"`
	Bridge  *Bridge  `xml:"bridge,omitempty"`
	Forward *Forward `xml:"forward,omitempty"`
	IP      IP       `xml:"ip"`
}

// Bridge represents the <bridge> element.
type Bridge struct {
	Name string `xml:"name,attr"`
}
type Forward struct {
	Mode string `xml:"mode,attr"`
}
type IP struct {
	Address string `xml:"address,attr"`
	Netmask string `xml:"netmask,attr"`
	// DHCP is omitted if nil.
	DHCP *DHCP `xml:"dhcp,omitempty"`
}

type DHCP struct {
	Range Range `xml:"range"`
}

type Range struct {
	Start string `xml:"start,attr"`
	End   string `xml:"end,attr"`
}

type NetworkOption func(*Network)

// WithBridge is an option that sets a custom bridge name.
// If an empty string is passed, it leaves the Bridge nil, allowing libvirt to create
// it automatically.
func WithBridge(bridgeName string) NetworkOption {
	return func(n *Network) {
		n.Bridge = &Bridge{Name: bridgeName}
	}
}

// WithDHCP is an option that enables DHCP if the flag is true, and sets the start and end range.
// If disabled (false), the DHCP field is left nil and omitted from the XML.

func WithDHCP(start, end string) NetworkOption {
	return func(n *Network) {
		n.IP.DHCP = &DHCP{
			Range: Range{
				Start: start,
				End:   end,
			},
		}
	}
}

// NewNetwork is the constructor that creates a new Network instance.
// it takes required parameters: name, forwardMode, ipAddress, and netmask.
// Addtional optional configurations can be provided using variadic options.

func NewNetwork(
	name, forwardMode, ipAddress, netmask string,
	autostart bool,
	opts ...NetworkOption,
) *Network {
	// Create the network with required fields.
	network := &Network{
		Name: name,
		Forward: &Forward{
			Mode: forwardMode,
		},
		// DHCP is nil by default, meaning it will be omitted unless enabled.
		IP: IP{
			Address: ipAddress,
			Netmask: netmask,
		},
	}
	for _, opt := range opts {
		opt(network)
	}
	return network
}

// GenerateXML returns the XML representation of the Network.
func (n *Network) GenerateXML() ([]byte, error) {
	return xml.MarshalIndent(n, "", "  ")
}
