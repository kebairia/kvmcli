package vms

import (
	"fmt"
	"path/filepath"
)

// Create Virtual Machine
func (vm *VirtualMachine) Create() error {
	// Initiliaze a new vm record
	record, err := NewVirtualMachineRecord(vm)
	if err != nil {
		return fmt.Errorf("can't Initiliaze a new vm: %w", err)
	}

	store, img, err := vm.fetchStoreAndImage(vm.Config.Spec.Image)
	if err != nil {
		return fmt.Errorf("fetch store and image: %w", err)
	}

	src := filepath.Join(store.ArtifactsPath, img.File)
	dest := filepath.Join(store.ImagesPath, vm.Config.Metadata.Name+".qcow2")

	// this is a slice that collection all the cleanup functions
	// so that when an error happens in each step, a proper rollback
	// will be executed
	var cleanups []func() error

	// artifactsPath, imagesPath := vm.disk.Paths()
	// src := fmt.Sprintf("%s/%s", artifactsPath, vm.Config.Spec.Image)
	// dest := fmt.Sprintf("%s/%s.qcow2", imagesPath, vm.Config.Metadata.Name)
	if err := vm.disk.CreateOverlay(vm.ctx, src, dest); err != nil {
		return fmt.Errorf("create disk overlay: %w", err)
	}
	//
	cleanups = append(cleanups, func() error {
		return vm.disk.DeleteOverlay(vm.ctx, dest)
	})

	// Generate the libvirt XML configuration
	xmlConfig, err := vm.domain.BuildXML(vm.ctx, vm.db, vm.Config)
	if err != nil {
		return vm.rollback(cleanups, "build XML", err)
	}

	// Step 3: Define and start the VM
	if err := vm.domain.Define(vm.ctx, xmlConfig); err != nil {
		return vm.rollback(cleanups, "define domain", err)
	}

	cleanups = append(cleanups, func() error {
		return vm.domain.Undefine(vm.ctx, vm.Config.Metadata.Name)
	})

	if err := vm.domain.Start(vm.ctx, vm.Config.Metadata.Name); err != nil {
		return vm.rollback(cleanups, "start domain", err)
	}

	// cleanups = append(cleanups, vm.domain.)
	cleanups = append(cleanups, func() error {
		return vm.domain.Stop(vm.ctx, vm.Config.Metadata.Name)
	})

	// Insert the vm record
	if err = record.Insert(vm.ctx, vm.db); err != nil {
		return vm.rollback(cleanups, "insert record", err)
	}

	fmt.Printf("vm/%s created\n", vm.Config.Metadata.Name)
	return nil
}
