package vms

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kebairia/kvmcli/internal/config"
	"github.com/kebairia/kvmcli/internal/logger"
	op "github.com/kebairia/kvmcli/internal/operations"
)

const imagesPath = "/home/zakaria/dox/homelab/images/"

func DestroyFromFile(configPath string) error {
	conn, err := op.InitConnection("unix", "/var/run/libvirt/libvirt-sock")
	if err != nil {
		return fmt.Errorf("failed to establish libvirt connection: %w", err)
	}
	defer conn.Disconnect()

	manager := NewVMManager(conn)

	var vms []config.VirtualMachine
	if vms, err = config.LoadConfig(configPath); err != nil {
		logger.Log.Errorf("failed to connect: %v", err)
	}

	for _, vm := range vms {
		if err := manager.Delete(vm.Metadata.Name); err != nil {
			logger.Log.Errorf("%s", err)
		}
	}
	return nil
}

func DestroyFromArgs(vmNames []string) error {
	conn, err := op.InitConnection("unix", "/var/run/libvirt/libvirt-sock")
	if err != nil {
		return fmt.Errorf("failed to establish libvirt connection: %w", err)
	}
	defer conn.Disconnect()
	manager := NewVMManager(conn)
	for _, vmName := range vmNames {
		if err := manager.Delete(vmName); err != nil {
			logger.Log.Errorf("error destroying VM %q: %v", vmName, err)
		}
	}
	return nil
}

func (m *VMManager) Delete(name string) error {
	// Lookup the domain by name
	domain, err := m.Conn.DomainLookupByName(name)
	if err != nil {
		return fmt.Errorf("Failed to find VM %s: %w", name, err)
	}
	logger.Log.Debugf("Deleting VM: %q", name)

	// Destroy the VM if it's running
	// NOTE: work on this later, remove only vm not running

	if err := m.Conn.DomainDestroy(domain); err != nil {
		// It might be acceptable if the VM is not running; log a warning instead of failing immediately.
		return fmt.Errorf("failed to destroy VM %q (it might not be running): %w", name, err)
	}

	// Undefine the domain
	if err := m.Conn.DomainUndefine(domain); err != nil {
		return fmt.Errorf("failed to undefine VM %q: %w", name, err)
	}

	logger.Log.Debugf("%q has been successfully undefined", name)
	// Remove the disk of the virtual machine
	diskPath := fmt.Sprintf("%s.qcow2", filepath.Join(imagesPath, name))
	if err := os.Remove(diskPath); err != nil {
		return fmt.Errorf("failed to delete disk for VM %q: %w", name, err)
	}

	logger.Log.Infof("%s/%s deleted", "vm", name)
	return nil
}
