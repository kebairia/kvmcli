package network

import (
	"github.com/digitalocean/go-libvirt"
)

// Struct definition
type VirtualNetwork struct {
	// Conn to hold the libvirt connection
	Conn       *libvirt.Libvirt `yaml:"-"`
	ApiVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   NetMetadata      `yaml:"metadata"`
	Spec       NetSpec          `yaml:"spec"`
}
type NetMetadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}
type NetSpec struct {
	MacAddress string            `yaml:"macAddress"`
	Bridge     string            `yaml:"bridge"`
	Mode       string            `yaml:"mode"`
	NetAddress string            `yaml:"netAddress"`
	Netmask    string            `yaml:"netmask"`
	DHCP       map[string]string `yaml:"dhcp"`
	Autostart  bool              `yaml:"autostart"`
}

func (net *VirtualNetwork) SetConnection(conn *libvirt.Libvirt) {
	net.Conn = conn
}
