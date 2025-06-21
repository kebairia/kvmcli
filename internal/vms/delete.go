package vms

import (
	"fmt"
	// db "github.com/kebairia/kvmcli/internal/database"
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
	if vm.conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}

	vmName := vm.Config.Metadata.Name
	domain, err := vm.conn.DomainLookupByName(vmName)
	if err != nil {
		// return fmt.Errorf("Failed to find VM %s: %w", vmName, err)
		return err
	}

	// Attempt to destroy the domain.
	if err := vm.conn.DomainDestroy(domain); err != nil {
		return err
	}

	// Undefine the domain
	if err := vm.conn.DomainUndefine(domain); err != nil {
		return fmt.Errorf("failed to undefine VM %q: %w", vmName, err)
	}

	// Remove the disk associated with the VM.
	if err := vm.CleanupDisk(); err != nil {
		return err
	}

	record, err := NewVirtualMachineRecord(vm)
	if err != nil {
		return err
	}
	err = record.Delete(vm.ctx, vm.db)
	if err != nil {
		return err
	}
	fmt.Printf("vm/%s deleted\n", vmName)

	return nil
}
