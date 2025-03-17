package operations

import (
	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/config"
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/utils"
)

// ProvisionVMs reads the server configuration file, establishes a connection to libvirt,
// and provisions each virtual machine defined in the configuration.
// It iterates over each VM entry, creates the corresponding domain configuration,
// generates the XML, creates an overlay disk, and then defines and starts the VM.

func CreateVMFromConfig(configPath string) error {
	// Establish a connection to libvirt.
	// The network type "unix" and the socket path are specified; these can be made configurable.
	libvirtConn, err := InitConnection("unix", "/var/run/libvirt/libvirt-sock")
	if err != nil {
		logger.Log.Fatalf("Failed to establish libvirt connection: %v", err)
	}
	// Ensure that the libvirt connection is closed when the function exits.
	defer libvirtConn.Disconnect()

	// Load server configuration from the YAML file.
	// The configuration file path is hardcoded; consider reading it from environment variables or flags.
	var vms []config.VirtualMachine
	if err := config.LoadConfig(configPath, &vms); err != nil {
		logger.Log.Fatal(err)
	}

	// serverConfig, err := config.LoadConfig[config.VirtualMachine](configPath)
	// if err != nil {
	// 	logger.Log.Fatal(err)
	// }

	// Iterate over the VMs defined in the configuration.
	for vmName, vmConfig := range serverConfig.VMs {
		logger.Log.Debugf("Provisioning VM: %s", vmName)

		// Create a domain definition from the VM configuration.
		// The NewDomain helper function constructs a domain object with proper settings.
		domain := utils.NewDomain(
			vmName,
			vmConfig.Memory,
			vmConfig.CPU,
			vmConfig.Disk.Path,
			vmConfig.Network.MAC,
		)
		// Create an overlay disk image based on a base image.
		// TODO: Ensure that CreateOverlay returns an error, and handle it appropriately.
		if err := CreateOverlay("rocky.qcow2", vmConfig.Disk.Path); err != nil {
			logger.Log.Errorf("Failed to create overlay for VM %s: %v", vmName, err)
		}

		// Generate the XML configuration required by libvirt for the VM.
		xmlConfig, err := domain.GenerateXML()
		if err != nil {
			logger.Log.Warnf("Failed to generate XML for VM %s: %v", vmName, err)
			continue
		}
		// Create the VM using the generated XML configuration.
		CreateVM(vmName, xmlConfig, libvirtConn)

	}
}

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

	logger.Log.Infof("%s/%s created", "vm", vmInstance.Name)
}
