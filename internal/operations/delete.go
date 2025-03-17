package operations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kebairia/kvmcli/internal/config"
	"github.com/kebairia/kvmcli/internal/logger"
)

const imagesPath = "/home/zakaria/dox/homelab/images/"

func DestroyFromFile(configPath string) {
	var vms []config.VirtualMachine
	var err error
	if vms, err = config.LoadConfig(configPath); err != nil {
		logger.Log.Errorf("failed to connect: %v", err)
	}

	for _, vm := range vms {
		if err := DestroyVM(vm.Metadata.Name); err != nil {
			logger.Log.Errorf("%s", err)
		}
	}
}

func DestroyFromArgs(vmNames []string) error {
	for _, vmName := range vmNames {
		if err := DestroyVM(vmName); err != nil {
			logger.Log.Errorf("error destroying VM %q: %v", vmName, err)
		}
	}
	return nil
}

func DestroyVM(vmName string) error {
	libvirtConn, err := InitConnection("unix", "/var/run/libvirt/libvirt-sock")
	if err != nil {
		return fmt.Errorf("failed to establish libvirt connection: %w", err)
	}
	defer libvirtConn.Disconnect()

	// Lookup the domain by name
	domain, err := libvirtConn.DomainLookupByName(vmName)
	if err != nil {
		return fmt.Errorf("Failed to find VM %s: %w", vmName, err)
	}
	logger.Log.Debugf("Deleting VM: %q", vmName)

	// Destroy the VM if it's running
	// NOTE: work on this later, remove only vm not running

	if err := libvirtConn.DomainDestroy(domain); err != nil {
		// It might be acceptable if the VM is not running; log a warning instead of failing immediately.
		return fmt.Errorf("failed to destroy VM %q (it might not be running): %w", vmName, err)
	}

	// Undefine the domain
	if err := libvirtConn.DomainUndefine(domain); err != nil {
		return fmt.Errorf("failed to undefine VM %q: %w", vmName, err)
	}

	logger.Log.Debugf("%q has been successfully undefined", vmName)
	// Remove the disk of the virtual machine
	diskPath := fmt.Sprintf("%s.qcow2", filepath.Join(imagesPath, vmName))
	if err := os.Remove(diskPath); err != nil {
		return fmt.Errorf("failed to delete disk for VM %q: %w", vmName, err)
	}

	logger.Log.Infof("%s/%s deleted", "vm", vmName)
	return nil
}
