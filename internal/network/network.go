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

var ErrVirtualNetworkNameEmpty = errors.New("virtual network name is empty")

// VirtualNetwork manages a libvirt-backed network, with IP⇔MAC mapping stored in DB.
type VirtualNetwork struct {
	Spec Network
	// Config VirtualNetworkConfig
	ctx  context.Context
	db   *sql.DB
	conn *libvirt.Libvirt
}

// VirtualNetworkOption configures a VirtualNetwork.
type VirtualNetworkOption func(*VirtualNetwork)

// WithLibvirtConnection sets the libvirt client (required).
func WithLibvirtConnection(conn *libvirt.Libvirt) VirtualNetworkOption {
	return func(vn *VirtualNetwork) {
		vn.conn = conn
	}
}

// WithDatabaseConnection sets the SQL database (required).
func WithDatabaseConnection(db *sql.DB) VirtualNetworkOption {
	return func(vn *VirtualNetwork) {
		vn.db = db
	}
}

// WithContext sets a custom context. If nil is passed, context.Background() is used.
func WithContext(ctx context.Context) VirtualNetworkOption {
	return func(vn *VirtualNetwork) {
		if ctx == nil {
			vn.ctx = context.Background()
		} else {
			vn.ctx = ctx
		}
	}
}

// NewVirtualNetwork creates a VirtualNetwork, applying options and validating dependencies.
func NewVirtualNetwork(
	// cfg VirtualNetworkConfig,
	spec Network,
	opts ...VirtualNetworkOption,
) (*VirtualNetwork, error) {
	// if cfg.Metadata.Name == "" {
	// 	return nil, ErrVirtualNetworkNameEmpty
	// }
	if spec.Name == "" {
		return nil, ErrVirtualNetworkNameEmpty
	}

	vn := &VirtualNetwork{
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
func (vn *VirtualNetwork) AddStaticMapping(ip, mac string) error {
	// TODO: load existing network XML via vn.conn.LookupNetworkByName
	// TODO: inject <host ip="..." mac="..."/> into XML
	// TODO: define vn.conn.NetworkDefineXML and vn.conn.NetworkUpdate call
	// TODO: persist mapping in vn.db
	return nil
}

// NOTE: this is just to change later
func (vn *VirtualNetwork) Start() error {
	fmt.Println("Start virtual network")
	return nil
}
