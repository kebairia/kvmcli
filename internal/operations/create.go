package operations

import (
	"fmt"
	"log"

	"github.com/digitalocean/go-libvirt"
)

// CreateVM creates a virtual machine using the provided VM name and XML configuration.
// It defines the domain in libvirt and then starts the domain.
// If an error occurs during either step, it logs the error.
// TODO: Consider returning an error instead of logging it directly for better error propagation.
func CreateVM(vmName string, xmlConfig []byte, conn *libvirt.Libvirt) {
	// Define the domain in libvirt using the provided XML configuration.
	vmInstance, err := conn.DomainDefineXML(string(xmlConfig))
	if err != nil {
		fmt.Printf("Error: Failed to define domain for VM %s: %v", vmName, err)
	}

	log.Printf("VM defined successfully: %s", vmInstance.Name)

	// Start the VM using the defined domain.
	if err := conn.DomainCreate(vmInstance); err != nil {
		fmt.Printf("Error: Failed to start VM %s: %v", vmInstance.Name, err)
	}

	log.Printf("VM started successfully: %s", vmInstance.Name)
}
