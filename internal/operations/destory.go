package operations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kebairia/kvmcli/internal/logger"
)

const imagesPath = "/home/zakaria/dox/homelab/images/"

func DestroyFromFile(path string) {
}

func DestroyFromArgs(...[]string) {
}

func DestroyVM(vmName string) {
	libvirtConn, err := InitConnection("unix", "/var/run/libvirt/libvirt-sock")
	if err != nil {
		logger.Log.Errorf("failed to establish libvirt connection: %v", err)
	}
	defer libvirtConn.Disconnect()

	// Lookup the domain by name
	domain, err := libvirtConn.DomainLookupByName(vmName)
	if err != nil {
		logger.Log.Errorf("Failed to find VM %s: %v", vmName, err)
	}
	logger.Log.Debugf("Deleting VM: %q", vmName)

	// Destroy the VM if it's running
	// NOTE: work on this later, remove only vm not running

	// state, someting, err := libvirtConn.DomainGetState(domain, 1)
	// fmt.Println(state, someting)

	if err := libvirtConn.DomainDestroy(domain); err != nil {
		// It might be acceptable if the VM is not running; log a warning instead of failing immediately.
		logger.Log.Fatalf("failed to destroy VM %q (it might not be running): %v", vmName, err)
	}

	// Undefine the domain
	if err := libvirtConn.DomainUndefine(domain); err != nil {
		logger.Log.Fatalf("failed to undefine VM %q: %v", vmName, err)
	}
	logger.Log.Debugf("%q has been successfully undefined", vmName)
	// Remove the disk of the virtual machine
	diskPath := fmt.Sprintf("%s.qcow2", filepath.Join(imagesPath, vmName))
	if err := os.Remove(diskPath); err != nil {
		logger.Log.Fatalf("failed to delete disk for VM %q: %v", vmName, err)
	}

	logger.Log.Infof("%q has been successfully deleted", vmName)
}
