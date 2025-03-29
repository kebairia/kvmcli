package vms

import (
	"fmt"

	db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
)

// Create Virtual Machine
func (vm *VirtualMachine) Create() error {
	// Check connection
	if vm.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}
	// Initiliaze a new vm record
	record := NewVMRecord(vm)

	// Insert the vm record
	_, err := db.Insert(record)
	if err != nil {
		return fmt.Errorf("failed to create database record for %q: %w", vm.Metadata.Name, err)
	}

	// Create overlay image
	// FIX: fix  error handling, add the CreateOverlay the error message below
	// for better context, and use logger.Log.Errorf here.

	if err := CreateOverlay("rocky-base-image.qcow2", vm.Spec.Disk.Path); err != nil {
		return fmt.Errorf("Failed to create overlay for VM %q: %w", vm.Metadata.Name, err)
	}

	// Prepare the domain and generate its XML configuration.
	xmlConfig, err := vm.prepareDomain()
	if err != nil {
		logger.Log.Errorf("%v", err)
	}
	// Define the domain and start the VM.
	if err := vm.defineAndStartDomain(xmlConfig); err != nil {
		logger.Log.Errorf("%v", err)
	}

	fmt.Printf("vm/%s created\n", vm.Metadata.Name)
	return nil
}
