package vms

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/digitalocean/go-libvirt"
	db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/utils"
)

const (
	deviceName = "vda"

	// Domain state constants.
	domainStateRunning = 1
	domainStatePaused  = 3
	domainStateStopped = 5
)

// CreateOverlay creates a qcow2 overlay image using a backing file obtained from the store record.
// It invokes the 'qemu-img' utility with a timeout context.
func (vm *VirtualMachine) CreateOverlay(image string) error {
	var st db.StoreRecord
	if err := st.GetRecord(db.Ctx, db.DB, "homelab-store"); err != nil {
		return fmt.Errorf("can't get store %q: %w", "homelab-store", err)
	}

	// Pull the image entry the VM asked for (imageKey could be vm.Spec.Image)
	img, ok := st.Images[image]
	if !ok {
		return fmt.Errorf("image %q not found in store", image)
	}

	// Build the full path to the base image from the store configuration.
	baseImagePath := filepath.Join(
		st.ArtifactsPath,
		img.Directory,
		img.File,
	)
	// Construct target overlay image file name.
	imageFile := fmt.Sprintf("%s.qcow2", filepath.Join(st.ImagesPath, vm.Metadata.Name))

	// Define the qemu-img command arguments.
	cmdArgs := []string{
		"create",
		"-o", fmt.Sprintf("backing_file=%s,backing_fmt=qcow2", baseImagePath),
		"-f", "qcow2",
		imageFile,
	}

	// Set a timeout context for running the external command.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute the command.
	cmd := exec.CommandContext(ctx, "qemu-img", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Debugf("qemu-img output: %s", output)
		return fmt.Errorf("failed to execute qemu-img command: %w", err)
	}

	logger.Log.Debug("Overlay image created successfully")
	return nil
}

// DeleteOverlay deletes the qcow2 overlay image file from the file system.
// It gets the target disk image path based on the store configuration and VM metadata.
func (vm *VirtualMachine) DeleteOverlay(image string) error {
	var st db.StoreRecord
	if err := st.GetRecord(db.Ctx, db.DB, "homelab-store"); err != nil {
		return fmt.Errorf("can't get store %q: %w", "homelab-store", err)
	}
	// Construct the disk image path.
	diskPath := filepath.Join(st.ImagesPath, vm.Metadata.Name+".qcow2")
	if err := os.Remove(diskPath); err != nil {
		return fmt.Errorf("failed to delete disk for VM %q: %w", vm.Metadata.Name, err)
	}
	return nil
}

func getNetworkIDByName(ctx context.Context, db *sql.DB, networkName string) (int, error) {
	const query = `
		SELECT id FROM networks
		WHERE name = ? 
	`

	var networkID int
	err := db.QueryRowContext(ctx, query, networkName).Scan(&networkID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve network ID for Network %q: %w", networkName, err)
	}

	return networkID, nil
}

// NewVMRecord constructs a new VM record from the provided VM configuration.
// The record is then used to insert VM metadata into the database.
func NewVMRecord(
	ctx context.Context,
	mydb *sql.DB,
	vm *VirtualMachine,
) (*db.VirtualMachineRecord, error) {
	var st db.StoreRecord
	if err := st.GetRecord(db.Ctx, db.DB, "homelab-store"); err != nil {
		return &db.VirtualMachineRecord{}, fmt.Errorf(
			"can't get store %q: %w",
			"homelab-store",
			err,
		)
	}

	// Build the disk image path (with a .qcow2 extension) based on the store configuration.
	diskImagePath := fmt.Sprintf(
		"%s.qcow2",
		filepath.Join(st.ImagesPath, vm.Metadata.Name),
	)
	// Lookup the network ID based on the network name
	networkID, err := getNetworkIDByName(ctx, mydb, vm.Spec.Network.Name)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve network ID: %w", err)
	}
	return &db.VirtualMachineRecord{
		Name:       vm.Metadata.Name,
		Namespace:  vm.Metadata.Namespace,
		Labels:     vm.Metadata.Labels,
		CPU:        vm.Spec.CPU,
		RAM:        vm.Spec.Memory,
		DiskSize:   vm.Spec.Disk.Size,
		DiskPath:   diskImagePath,
		Image:      vm.Spec.Image,
		MacAddress: vm.Spec.Network.MacAddress,
		NetworkID:  networkID,
		CreatedAt:  time.Now(),
	}, nil
}

// prepareDomain generates the XML configuration for the virtual machine domain.
// It uses the store record to determine the disk image location and creates the domain configuration.
func (vm *VirtualMachine) prepareDomain(image string) (string, error) {
	// Build the full path to the disk image with the .qcow2 extension.
	var st db.StoreRecord
	if err := st.GetRecord(db.Ctx, db.DB, "homelab-store"); err != nil {
		return "", fmt.Errorf("can't get store %q: %w", "homelab-store", err)
	}
	// Pull the image entry the VM asked for (imageKey could be vm.Spec.Image)
	img, ok := st.Images[image]
	if !ok {
		return "", fmt.Errorf("image %q not found in store", image)
	}
	// Build the disk image path for the domain configuration.
	diskImagePath := fmt.Sprintf(
		"%s.qcow2",
		filepath.Join(st.ImagesPath, vm.Metadata.Name),
	)

	// Create a new domain configuration using utility functions.
	domain := utils.NewDomain(
		vm.Metadata.Name,
		vm.Spec.Memory,
		vm.Spec.CPU,
		diskImagePath,
		vm.Spec.Network.Name,
		vm.Spec.Network.MacAddress,
		img.OsProfile,
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

// getState returns a string representation of the VM state based on its domain info.
func getState(conn *libvirt.Libvirt, domain libvirt.Domain) (string, error) {
	state, _, _, _, _, err := conn.DomainGetInfo(domain)
	if err != nil {
		return "", fmt.Errorf("failed to get info for domain %s: %w", domain.Name, err)
	}

	switch int(state) {
	case domainStateRunning:
		return "Running", nil
	case domainStatePaused:
		return "Paused", nil
	case domainStateStopped:
		return "Stopped", nil
	default:
		return "Unknown", nil
	}
}

// getDiskSize returns the disk size (in gigabytes) for the specified VM domain.
func getDiskSize(conn *libvirt.Libvirt, domain libvirt.Domain) (float64, error) {
	_, _, diskPhysSize, err := conn.DomainGetBlockInfo(domain, deviceName, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to get block info for domain %s: %w", domain.Name, err)
	}

	return float64(diskPhysSize) / (1024 * 1024 * 1024), nil
}

// formatAge returns a human-friendly string for the time elapsed since t.
func formatAge(t time.Time) string {
	duration := time.Since(t)
	if duration < 0 {
		duration = -duration
	}

	if days := int(duration.Hours() / 24); days >= 1 {
		return fmt.Sprintf("%dd", days)
	}
	if hours := int(duration.Hours()); hours >= 1 {
		return fmt.Sprintf("%dh", hours)
	}
	if minutes := int(duration.Minutes()); minutes >= 1 {
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%ds", int(duration.Seconds()))
}
