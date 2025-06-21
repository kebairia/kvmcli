package vms

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/digitalocean/go-libvirt"
	db "github.com/kebairia/kvmcli/internal/database"
	log "github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/utils"
)

const (
	deviceName = "vda"
	// qemu
	defaultQemuTimeout = 30 * time.Second
	qemuImgCmd         = "qemu-img"
)

// SetupDisk creates a QCOW2 overlay image for the VM's disk.
func (vm *VirtualMachine) SetupDisk() error {
	store, img, err := vm.fetchStoreAndImage(vm.Config.Spec.Image)
	if err != nil {
		return fmt.Errorf("fetch store and image: %w", err)
	}

	src := filepath.Join(store.ArtifactsPath, img.File)
	dest := filepath.Join(store.ImagesPath, vm.Config.Metadata.Name+".qcow2")

	if err := vm.disk.CreateOverlay(vm.ctx, src, dest); err != nil {
		return fmt.Errorf("create overlay image: %w", err)
	}

	return nil
}

// CleanupDisk removes the VM's QCOW2 overlay file from disk.
func (vm *VirtualMachine) CleanupDisk() error {
	store, err := vm.fetchStore()
	if err != nil {
		return fmt.Errorf("fetch store: %w", err)
	}

	// Ensure the image exists in the store
	if _, err := store.GetImageRecord(vm.ctx, vm.db, vm.Config.Spec.Image); err != nil {
		return fmt.Errorf("image %q not found in store: %w", vm.Config.Spec.Image, err)
	}

	diskPath := filepath.Join(store.ImagesPath, vm.Config.Metadata.Name+".qcow2")
	if err := os.Remove(diskPath); err != nil {
		return fmt.Errorf("remove disk %q: %w", diskPath, err)
	}

	// dest := filepath.Join(store.ImagesPath, vm.Config.Metadata.Name+".qcow2")
	//
	// if err := vm.disk.DeleteOverlay(vm.ctx, dest); err != nil {
	// 	return err
	// }

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
// then creates a database record for the virtual machine.
func NewVirtualMachineRecord(
	vm *VirtualMachine,
) (*db.VirtualMachineRecord, error) {
	store, err := vm.fetchStore()
	if err != nil {
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	// Verify image exists in store
	if _, err := store.GetImageRecord(vm.ctx, vm.db, vm.Config.Spec.Image); err != nil {
		return nil, fmt.Errorf(
			"image %q not found in store %q: %w",
			vm.Config.Spec.Image,
			vm.Config.Metadata.Store,
			err,
		)
	}

	networkID, err := getNetworkIDByName(vm.ctx, vm.db, vm.Config.Spec.Network.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}

	storeID, err := db.GetStoreIDByName(vm.ctx, vm.db, vm.Config.Metadata.Store)
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}
	diskPath := filepath.Join(store.ImagesPath, vm.Config.Metadata.Name+".qcow2")

	return &db.VirtualMachineRecord{
		Name:       vm.Config.Metadata.Name,
		Namespace:  vm.Config.Metadata.Namespace,
		Labels:     vm.Config.Metadata.Labels,
		CPU:        vm.Config.Spec.CPU,
		RAM:        vm.Config.Spec.Memory,
		DiskSize:   vm.Config.Spec.Disk.Size,
		DiskPath:   diskPath,
		Image:      vm.Config.Spec.Image,
		MacAddress: vm.Config.Spec.Network.MacAddress,
		IP:         vm.Config.Spec.Network.IP,
		NetworkID:  networkID,
		StoreID:    storeID,
		CreatedAt:  time.Now(),
	}, nil
}

// prepareDomain generates the XML configuration for the virtual machine domain.
// It uses the store record to determine the disk image location and creates the domain configuration.
func (vm *VirtualMachine) buildDomainXML() (string, error) {
	// Build the full path to the disk image with the .qcow2 extension.
	var st db.StoreRecord
	var err error
	st.ID, err = db.GetStoreIDByName(vm.ctx, vm.db, vm.Config.Metadata.Store)
	img, err := st.GetImageRecord(vm.ctx, vm.db, vm.Config.Spec.Image)
	if err != nil {
		return "", fmt.Errorf("can't get store %q: %w", vm.Config.Metadata.Store, err)
	}
	// Build the disk image path for the domain configuration.
	diskImagePath := fmt.Sprintf(
		"%s.qcow2",
		filepath.Join(st.ImagesPath, vm.Config.Metadata.Name),
	)

	// Create a new domain configuration using utility functions.
	domain := utils.NewDomain(
		vm.Config.Metadata.Name,
		vm.Config.Spec.Memory,
		vm.Config.Spec.CPU,
		diskImagePath,
		vm.Config.Spec.Network.Name,
		vm.Config.Spec.Network.MacAddress,
		img.OsProfile,
	)

	xmlConfig, err := domain.GenerateXML()
	if err != nil {
		return "", fmt.Errorf("failed to generate XML for VM %s: %v", vm.Config.Metadata.Name, err)
	}

	return xml.Header + string(xmlConfig), nil
}

// defineAndStartDomain defines the domain using the provided XML configuration and starts the VM.
func (vm *VirtualMachine) defineAndStartDomain(xmlConfig string) error {
	vmInstance, err := vm.conn.DomainDefineXML(xmlConfig)
	if err != nil {
		return fmt.Errorf("failed to define domain for VM %s: %v", vm.Config.Metadata.Name, err)
	}

	if err := vm.conn.DomainCreate(vmInstance); err != nil {
		return fmt.Errorf("failed to start VM %s: %w", vm.Config.Metadata.Name, err)
	}

	return nil
}

func (vm *VirtualMachine) undefineAndDestoryDomain() error {
	domain, err := vm.conn.DomainLookupByName(vm.Config.Metadata.Name)
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
		return fmt.Errorf("failed to undefine VM %q: %w", vm.Config.Metadata.Name, err)
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

// getStore retrieves the store record for the VM.
func (vm *VirtualMachine) fetchStore() (*db.StoreRecord, error) {
	var store db.StoreRecord
	var err error

	store.ID, err = db.GetStoreIDByName(vm.ctx, vm.db, vm.Config.Metadata.Store)
	if err != nil {
		return nil, fmt.Errorf("failed to get store ID for %q: %w", vm.Config.Metadata.Store, err)
	}

	return &store, nil
}

// getStoreAndImage retrieves both store and image records.
func (vm *VirtualMachine) fetchStoreAndImage(
	imageName string,
) (*db.StoreRecord, *db.ImageRecord, error) {
	store, err := vm.fetchStore()
	if err != nil {
		return nil, nil, err
	}

	img, err := store.GetImageRecord(vm.ctx, vm.db, imageName)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get image %q from store %q: %w",
			imageName,
			vm.Config.Metadata.Store,
			err,
		)
	}

	return store, img, nil
}

// // buildOverlayPath constructs the full path for the overlay image.
// func (vm *VirtualMachine) buildOverlayPath(store *db.StoreRecord) string {
// 	return filepath.Join(store.ImagesPath, vm.Metadata.Name+".qcow2")
// }

func (vm *VirtualMachine) rollback(cleanups []func() error, step string, originError error) error {
	for _, fn := range cleanups {
		if err := fn(); err != nil {
			log.Warnf("rollback failed, step %s, err %s", step, err)
		}
	}

	return fmt.Errorf("failed at %s: %w", step, originError)
}
