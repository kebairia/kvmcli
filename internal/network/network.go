package network

import (
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/logger"
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

func (net *VirtualNetwork) Create() error {
	// Check connection
	if net.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}
	xmlConfig, err := net.prepareNetwork()
	if err != nil {
		logger.Log.Fatalf("%v", err)
	}
	// Define the network and start it
	if err := net.defineAndStartNetwork(xmlConfig); err != nil {
		logger.Log.Errorf("%v", err)
	}

	fmt.Printf("net/%s created\n", net.Metadata.Name)
	return nil
}

func (net *VirtualNetwork) Delete() error {
	fmt.Println("net get deleted")
	return nil
}

func (net *VirtualNetwork) SetConnection(conn *libvirt.Libvirt) {
	net.Conn = conn
}
