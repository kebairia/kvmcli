package vms

import (
	"fmt"

	db "github.com/kebairia/kvmcli/internal/database"
)

// Create Virtual Machine
func (vm *VirtualMachine) Create() error {
	// Check connection
	if vm.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}
	// Initiliaze a new vm record
	record := NewVMRecord(vm)

	// Step 1: Create the overlay disk image
	if err := vm.CreateOverlay(vm.Spec.Image); err != nil {
		return fmt.Errorf("failed to create overlay for VM %q: %w", vm.Metadata.Name, err)
	}

	// Step 2: Generate the libvirt XML configuration
	xmlConfig, err := vm.prepareDomain()
	if err != nil {
		_ = vm.DeleteOverlay(vm.Metadata.Name) // rollback overlay image
		return fmt.Errorf("failed to prepare domain for VM %q: %w", vm.Metadata.Name, err)
	}

	// Step 3: Define and start the VM
	if err := vm.defineAndStartDomain(xmlConfig); err != nil {
		_ = vm.DeleteOverlay(
			vm.Metadata.Name,
		) // rollback, delete overlay if the domain preparation failed
		return fmt.Errorf("failed to define/start VM %q: %w", vm.Metadata.Name, err)
	}

	// Step 4: Insert the vm record
	if _, err = db.InsertVM(record); err != nil {
		_ = vm.Delete() // rollback libvirt domain and disk
		return fmt.Errorf("failed to create database record for VM %q: %w", vm.Metadata.Name, err)
	}

	fmt.Printf("vm/%s created\n", vm.Metadata.Name)
	return nil
}
