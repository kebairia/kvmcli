package network

import (
	"context"
	"database/sql"
	"errors"

	"github.com/digitalocean/go-libvirt"
)

var VirtualNetworkNameEmpty = errors.New("virtual network name is empty")

// Struct definition
type VirtualNetwork struct {
	// Conn to hold the libvirt connection
	Conn       *libvirt.Libvirt `yaml:"-"`
	DB         *sql.DB          `yaml:"-"`
	Context    context.Context  `yaml:"_"`
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

func (net *VirtualNetwork) SetConnection(ctx context.Context, db *sql.DB, conn *libvirt.Libvirt) {
	net.Conn = conn
	net.DB = db
	net.Context = ctx
}
