package network

import (
	"context"
)

// Start starts the network.
func (m *LibvirtNetworkManager) Start(ctx context.Context, name string) error {
	// TODO: Implement start logic if separate from creation.
	// For now, Create also Starts.
	// If Start is called separately, we should lookup and create.
	return nil
}

// AddStaticMapping adds a static IP mapping to the network.
func (m *LibvirtNetworkManager) AddStaticMapping(ctx context.Context, name, ip, mac string) error {
	// TODO: load existing network XML via m.conn.LookupNetworkByName
	// TODO: inject <host ip="..." mac="..."/> into XML
	// TODO: define m.conn.NetworkDefineXML and m.conn.NetworkUpdate call
	// TODO: persist mapping in m.db
	return nil
}
