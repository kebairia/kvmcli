package network

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/logger"
)

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
