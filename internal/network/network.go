package network

import (
	"fmt"
	"os"
	"text/tabwriter"

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
	DHCP       map[string]string `yaml:"dhcp"`
	Bridge     string            `yaml:"bridge"`
	Mode       string            `yaml:"mode"`
	Network    Network           `yaml:"network"`
	Autostart  bool              `yaml:"autostart"`
	MacAddress string            `yaml:"macAddress"`
}

type Network struct {
	Address string `yaml:"address"`
	Netmask string `yaml:"netmask"`
}

func (net *VirtualNetwork) SetConnection(conn *libvirt.Libvirt) {
	net.Conn = conn
}

func (net *VirtualNetwork) Header() *tabwriter.Writer {
	// Setup tabwriter for clean columnar output.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATE\tBRIDGE\tSUBNET\tGATEWAY\tDHCP RANGE\tAGE")

	return w
}

func (net *VirtualNetwork) PrintRow(w *tabwriter.Writer, info *VirtualNetworkInfo) {
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		info.Name,
		info.State,
		info.Bridge,
		info.Subnet,
		info.Gateway,
		info.DHCPRange,
		info.Age,
	)
}
