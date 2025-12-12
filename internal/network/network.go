package network

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/digitalocean/go-libvirt"
	// "github.com/kebairia/kvmcli/internal/config"
)

// IDEA: the ip <=>  mac address mapping done here in Virtual Network declaration
// What I need is whenever I create a new virtual machine with a static ip, I need to update
// my virtual network declaration to add the ip <=> mac address mapping

var ErrNetworkNameEmpty = errors.New("virtual network name is empty")

// Network manages a libvirt-backed network, with IP⇔MAC mapping stored in DB.
type Network struct {
	Spec Config
	// Config NetworkConfig
	ctx  context.Context
	db   *sql.DB
	conn *libvirt.Libvirt
}

// NetworkOption configures a Network.
type NetworkOption func(*Network)

// WithLibvirtConnection sets the libvirt client (required).
func WithLibvirtConnection(conn *libvirt.Libvirt) NetworkOption {
	return func(vn *Network) {
		vn.conn = conn
	}
}

// WithDatabaseConnection sets the SQL database (required).
func WithDatabaseConnection(db *sql.DB) NetworkOption {
	return func(vn *Network) {
		vn.db = db
	}
}

// WithContext sets a custom context. If nil is passed, context.Background() is used.
func WithContext(ctx context.Context) NetworkOption {
	return func(vn *Network) {
		if ctx == nil {
			vn.ctx = context.Background()
		} else {
			vn.ctx = ctx
		}
	}
}

// NewNetwork creates a Network, applying options and validating dependencies.
func NewNetwork(
	// cfg NetworkConfig,
	spec Config,
	opts ...NetworkOption,
) (*Network, error) {
	// if cfg.Metadata.Name == "" {
	// 	return nil, ErrNetworkNameEmpty
	// }
	if spec.Name == "" {
		return nil, ErrNetworkNameEmpty
	}

	vn := &Network{
		// Config: cfg,
		Spec: spec,
		// default context
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(vn)
	}

	// if vn.conn == nil {
	// 	return nil, ErrNilLibvirtConn
	// }
	// if vn.db == nil {
	// 	return nil, ErrNilDBConn
	// }

	return vn, nil
}

// AddStaticMapping records an IP⇔MAC mapping in the libvirt network XML and persists to DB.
func (vn *Network) AddStaticMapping(ip, mac string) error {
	// TODO: load existing network XML via vn.conn.LookupNetworkByName
	// TODO: inject <host ip="..." mac="..."/> into XML
	// TODO: define vn.conn.NetworkDefineXML and vn.conn.NetworkUpdate call
	// TODO: persist mapping in vn.db
	return nil
}

// NOTE: this is just to change later
func (vn *Network) Start() error {
	fmt.Println("Start virtual network")
	return nil
}
