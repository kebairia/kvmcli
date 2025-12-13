package network

import (
	"context"
	"database/sql"

	"github.com/digitalocean/go-libvirt"
)

// NetworkManager defines the interface for managing virtual networks.
type NetworkManager interface {
	Create(ctx context.Context, spec Config) error
	Delete(ctx context.Context, name string) error
	Start(ctx context.Context, name string) error
	AddStaticMapping(ctx context.Context, name, ip, mac string) error
}

// LibvirtNetworkManager implements NetworkManager using libvirt and a SQL database.
type LibvirtNetworkManager struct {
	conn *libvirt.Libvirt
	db   *sql.DB
}

// NewLibvirtNetworkManager creates a new LibvirtNetworkManager.
func NewLibvirtNetworkManager(conn *libvirt.Libvirt, db *sql.DB) *LibvirtNetworkManager {
	return &LibvirtNetworkManager{
		conn: conn,
		db:   db,
	}
}
