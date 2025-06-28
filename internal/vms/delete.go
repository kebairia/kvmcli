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

	vmName := vm.Config.Metadata.Name

	img, err := database.GetImageRecord(vm.ctx, vm.db, vm.Config.Spec.Image)
	if err != nil {
		return fmt.Errorf("fetch store and image: %w", err)
	}

	dest := filepath.Join(img.ImagesPath, vm.Config.Metadata.Name+".qcow2")

	// Attempt to destroy the domain.
	if err := vm.domain.Destroy(vm.ctx, vm.Config.Metadata.Name); err != nil {
		return err
	}

	// Undefine the domain
	if err := vm.domain.Undefine(vm.ctx, vm.Config.Metadata.Name); err != nil {
		return fmt.Errorf("failed to undefine VM %q: %w", vmName, err)
	}

	// Remove the disk associated with the VM.
	if err := vm.disk.DeleteOverlay(vm.ctx, dest); err != nil {
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
