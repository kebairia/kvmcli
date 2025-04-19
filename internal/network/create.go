package network

import (
	"fmt"

	databasesql "github.com/kebairia/kvmcli/internal/database-sql"
	"github.com/kebairia/kvmcli/internal/logger"
)

func (net *VirtualNetwork) Create() error {
	// Check connection
	if net.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}
	// record := NewNetRecord(net)
	record := NewVirtualNetworkRecord(net)

	// Insert the net record
	// _, err := db.InsertNet(record)
	err := databasesql.InsertNet(databasesql.Ctx, databasesql.DB, record)
	if err != nil {
		return fmt.Errorf("failed to create database record for %q: %w", net.Metadata.Name, err)
	}
	// Network
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
