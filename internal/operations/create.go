package operations

import (
	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/logger"
)

// CreateVM creates a virtual machine using the provided VM name and XML configuration.
// It defines the domain in libvirt and then starts the domain.
// If an error occurs during either step, it logs the error.
// TODO: Consider returning an error instead of logging it directly for better error propagation.
func CreateVM(vmName string, xmlConfig []byte, conn *libvirt.Libvirt) {
	// Define the domain in libvirt using the provided XML configuration.
	vmInstance, err := conn.DomainDefineXML(string(xmlConfig))
	if err != nil {
		logger.Log.Fatalf("Failed to define domain for VM %s: %v", vmName, err)
	}

	logger.Log.Debugf("%q defined successfully.", vmInstance.Name)

	// Start the VM using the defined domain.
	if err := conn.DomainCreate(vmInstance); err != nil {
		logger.Log.Fatalf("Failed to start VM %s: %v", vmInstance.Name, err)
	}

	logger.Log.Infof("%q started successfully", vmInstance.Name)
}
