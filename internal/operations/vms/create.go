package vms

import (
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/config"
	"github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
	op "github.com/kebairia/kvmcli/internal/operations"
	"github.com/kebairia/kvmcli/internal/utils"
)

// VMManager encapsulates the dependencies for managing VMs
type VMManager struct {
	Conn *libvirt.Libvirt
	// Logger *logger.Log
	// Log logrus.New()
}

func NewVMManager(conn *libvirt.Libvirt) *VMManager {
	return &VMManager{Conn: conn}
}

// ProvisionVMs reads the server configuration file, establishes a connection to libvirt,
// and provisions each virtual machine defined in the configuration.
// It iterates over each VM entry, creates the corresponding domain configuration,
// generates the XML, creates an overlay disk, and then defines and starts the VM.

func CreateVMFromConfig(configPath string) error {
	// Establish a connection to libvirt.
	// The network type "unix" and the socket path are specified; these can be made configurable.
	libvirtConn, err := op.InitConnection("unix", "/var/run/libvirt/libvirt-sock")
	if err != nil {
		logger.Log.Fatalf("Failed to establish libvirt connection: %v", err)
	}
	// Ensure that the libvirt connection is closed when the function exits.
	defer libvirtConn.Disconnect()

	// Load server configuration from the YAML file.
	vms, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Iterate over the VMs defined in the configuration.
	// Initilize a new manager for vm
	manager := NewVMManager(libvirtConn)
	for _, vm := range vms {
		logger.Log.Debugf("Provisioning VM: %s", vm.Metadata.Name)
		// Create a domain definition from the VM configuration.
		// The NewDomain helper function constructs a domain object with proper settings.
		// FIX: convert Memory into "1024MiB"
		// memoryStr := utils.FormatMemory(vm.Spec.Memory)

		domain := utils.NewDomain(
			vm.Metadata.Name,
			vm.Spec.Memory,
			vm.Spec.CPU,
			vm.Spec.Disk.Path,
			vm.Spec.Network.MacAddress,
		)

		// Create an overlay disk image based on a base image.
		if err := CreateOverlay("rocky.qcow2", vm.Spec.Disk.Path); err != nil {
			logger.Log.Errorf("Failed to create overlay for VM %s: %v", vm.Metadata.Name, err)
			// Continue to next VM even if overlay creation fails.
			continue
		}

		// Generate the XML configuration required by libvirt for the VM.
		xmlConfig, err := domain.GenerateXML()
		if err != nil {
			logger.Log.Warnf("Failed to generate XML for VM %s: %v", vm.Metadata.Name, err)
			continue
		}
		// Create a virtual machine using the provided name and xmlconfig file
		if err := manager.Create(vm.Metadata.Name, xmlConfig); err != nil {
			logger.Log.Errorf("%s", err)
		}

		// Create VM entry on database
		database.CreateVMEntry(
			vm.Metadata.Name,
			vm.Metadata.Namespace,
			vm.Spec.Memory,
			vm.Spec.CPU,
			vm.Spec.Network.MacAddress,
			vm.Spec.Network.Name,
		)
	}
	return nil
}

// CreateVM creates a virtual machine using the provided VM name and XML configuration.
// It defines the domain in libvirt and then starts the domain.
// If an error occurs during either step, it logs the error.

func (m *VMManager) Create(name string, xmlConfig []byte) error {
	// Define the domain in libvirt using the provided XML configuration.
	domain, err := m.Conn.DomainDefineXML(string(xmlConfig))
	if err != nil {
		return fmt.Errorf("Failed to define domain for VM %s: %w", name, err)
	}

	logger.Log.Debugf("%q defined successfully.", domain.Name)

	// Start the VM using the defined domain.
	if err := m.Conn.DomainCreate(domain); err != nil {
		return fmt.Errorf("Failed to start VM %s: %w", domain.Name, err)
	}

	logger.Log.Infof("%s/%s created", "vm", domain.Name)
	return nil
}
