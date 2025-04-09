package vms

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
)

// OPTIMIZE:

// 0. Check if connection is valide
// 1. Destroy and Undefine
// 2. Remove disk associated with the VM
// 3. Delete VM record from database

// A. DeleteMany for mongodb

// Delete Function
func (vm *VirtualMachine) Delete() error {
	var err error
	// Check connection
	// if connectionIsValide(vm.Conn), then (this logic)
	if vm.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}

	vmName := vm.Metadata.Name
	domain, err := vm.Conn.DomainLookupByName(vmName)
	if err != nil {
		return fmt.Errorf("Failed to find VM %s: %w", vmName, err)
	}

	// Attempt to destroy the domain.
	if err := vm.Conn.DomainDestroy(domain); err != nil {
		return fmt.Errorf(
			"failed to delete VM %q (it might not be running): %w",
			vmName,
			err,
		)
	}

	// Undefine the domain
	if err := vm.Conn.DomainUndefine(domain); err != nil {
		return fmt.Errorf("failed to undefine VM %q: %w", vmName, err)
	}

	// Remove the disk associated with the VM.
	if err := vm.DeleteOverlay(vmName); err != nil {
		return fmt.Errorf("failed to delete disk for VM %q: %w", vmName, err)
	}

	err = database.Delete(vm.Metadata.Name, database.VMsCollection)
	if err != nil {
		logger.Log.Errorf("failed to delete record for VM %s: %v", vm.Metadata.Name, err)
	}
	logger.Log.Infof("%s/%s deleted", "vm", vmName)

	return nil
}
