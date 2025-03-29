package vms

import (
	"context"
	"encoding/xml"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const artifactsPath = "/home/zakaria/dox/homelab/artifacts/rocky"

// CreateOverlay creates a qcow2 overlay image based on a backing file.
func CreateOverlay(baseImage, destImage string) error {
	// Construct the full path for the base image.
	baseImagePath := filepath.Join(artifactsPath, baseImage)

	// Prepare the qemu-img command.
	cmdArgs := []string{
		"create",
		"-o", fmt.Sprintf("backing_file=%s,backing_fmt=qcow2", baseImagePath),
		"-f", "qcow2",
		destImage,
	}

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "qemu-img", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Debugf("%s", output)
		return fmt.Errorf("failed execute qemu-img command: %w", err)
	}

	logger.Log.Debug("Overlay image created successfully")
	return nil
}

func NewVMRecord(vm *VirtualMachine) *db.VMRecord {
	// Create vm record out of infos
	return &db.VMRecord{
		Name:      vm.Metadata.Name,
		Namespace: vm.Metadata.Namespace,
		Labels:    vm.Metadata.Labels,
		CPU:       vm.Spec.CPU,
		RAM:       vm.Spec.Memory,
		Disk: db.Disk{
			vm.Spec.Disk.Size,
			vm.Spec.Disk.Path,
		},
		Image:       vm.Spec.Image,
		MacAddress:  vm.Spec.Network.MacAddress,
		Network:     vm.Spec.Network.Name,
		SnapshotIDs: []primitive.ObjectID{},
		CreatedAt:   time.Now(),
	}
}

func (vm *VirtualMachine) prepareDomain() (string, error) {
	// Creating domain out of infos
	domain := utils.NewDomain(

		vm.Metadata.Name,
		vm.Spec.Memory,
		vm.Spec.CPU,
		vm.Spec.Disk.Path,
		vm.Spec.Network.Name,
		vm.Spec.Network.MacAddress,
		"http://rockylinux.org/rocky/9",
	)

	xmlConfig, err := domain.GenerateXML()
	if err != nil {
		return "", fmt.Errorf("failed to generate XML for VM %s: %v", vm.Metadata.Name, err)
	}
	return fmt.Sprint(xml.Header + string(xmlConfig)), nil
}

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
