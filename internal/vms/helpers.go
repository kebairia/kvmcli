package vms

import (
	"context"
	"encoding/xml"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kebairia/kvmcli/internal/database"
	db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateOverlay creates a qcow2 overlay image based on a backing file.
func (vm *VirtualMachine) CreateOverlay(image string) error {
	st, err := database.GetRecord[database.StoreRecord](
		"local-image-store",
		database.StoreCollection,
	)
	if err != nil {
		return fmt.Errorf("can't get store %q: %w", "local-image-store", err)
	}

	baseImagePath := filepath.Join(
		st.Spec.Config.Path,
		st.Spec.Images[image].Directory,
		st.Spec.Images[image].File,
	)
	imageFile := fmt.Sprintf("%s.qcow2", filepath.Join(imagesPath, vm.Metadata.Name))

	cmdArgs := []string{
		"create",
		"-o", fmt.Sprintf("backing_file=%s,backing_fmt=qcow2", baseImagePath),
		"-f", "qcow2",
		imageFile,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "qemu-img", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Debugf("qemu-img output: %s", output)
		return fmt.Errorf("failed to execute qemu-img command: %w", err)
	}

	logger.Log.Debug("Overlay image created successfully")
	return nil
}

// NewVMRecord constructs a new VM record from the provided virtual machine information.
func NewVMRecord(vm *VirtualMachine) *db.VMRecord {
	diskImagePath := fmt.Sprintf("%s.qcow2", filepath.Join(imagesPath, vm.Metadata.Name))
	return &db.VMRecord{
		Name:      vm.Metadata.Name,
		Namespace: vm.Metadata.Namespace,
		Labels:    vm.Metadata.Labels,
		CPU:       vm.Spec.CPU,
		RAM:       vm.Spec.Memory,
		Disk: db.Disk{
			Size: vm.Spec.Disk.Size,
			Path: diskImagePath,
		},
		Image:       vm.Spec.Image,
		MacAddress:  vm.Spec.Network.MacAddress,
		Network:     vm.Spec.Network.Name,
		SnapshotIDs: []primitive.ObjectID{},
		CreatedAt:   time.Now(),
	}
}

// prepareDomain generates the XML configuration for the virtual machine domain.
func (vm *VirtualMachine) prepareDomain() (string, error) {
	// Build the full path to the disk image with the .qcow2 extension.
	diskImagePath := fmt.Sprintf("%s.qcow2", filepath.Join(imagesPath, vm.Metadata.Name))

	// Create a new domain configuration.
	domain := utils.NewDomain(
		vm.Metadata.Name,
		vm.Spec.Memory,
		vm.Spec.CPU,
		diskImagePath,
		vm.Spec.Network.Name,
		vm.Spec.Network.MacAddress,
		"http://rockylinux.org/rocky/9",
	)

	xmlConfig, err := domain.GenerateXML()
	if err != nil {
		return "", fmt.Errorf("failed to generate XML for VM %s: %v", vm.Metadata.Name, err)
	}

	// Prepend the XML header and return.
	return xml.Header + string(xmlConfig), nil
}

// defineAndStartDomain defines the domain using the provided XML configuration and starts the VM.
func (vm *VirtualMachine) defineAndStartDomain(xmlConfig string) error {
	vmInstance, err := vm.Conn.DomainDefineXML(xmlConfig)
	if err != nil {
		return fmt.Errorf("failed to define domain for VM %s: %v", vm.Metadata.Name, err)
	}

	if err := vm.Conn.DomainCreate(vmInstance); err != nil {
		return fmt.Errorf("failed to start VM %s: %w", vm.Metadata.Name, err)
	}

	return nil
}
