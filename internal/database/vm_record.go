package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	// "github.com/digitalocean/go-libvirt"
	// "github.com/kebairia/kvmcli/internal/logger"
	// "github.com/kebairia/kvmcli/internal/vms"
)

// VMRecord represents a VM stored in the SQLite database.
type VirtualMachineRecord struct {
	ID         int
	Name       string
	Namespace  string
	Labels     map[string]string
	CPU        int
	RAM        int
	MacAddress string
	NetworkID  int
	Image      string
	DiskSize   string
	DiskPath   string
	CreatedAt  time.Time
	// SnapshotIDs []string we don't use snapshot id here, in the snapshot table we reference  t the vm
}

// EnsureVMTable creates the vms table if it doesn't exist.
func EnsureVMTable(ctx context.Context, db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS vms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    namespace TEXT,
    cpu INTEGER,
    ram INTEGER,
    mac_address TEXT,
    network_id INTEGER,
    image TEXT,
    disk_size TEXT,
    disk_path TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    labels TEXT,
    FOREIGN KEY (network_id) REFERENCES networks(id)
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_vm_name_namespace ON vms(name, namespace);
	`
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create vms table: %w", err)
	}
	return nil
}

func (vmr *VirtualMachineRecord) GetRecord(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	query := fmt.Sprintf(`
		SELECT id, name, namespace, 
		       cpu, ram, mac_address, 
		       network_id, image, 
		       disk_size, disk_path, 
		       created_at, labels
		FROM %s
		WHERE name = ?`, VMsTable)

	// record    VirtualMachineRecord
	var labelText string

	err := db.QueryRowContext(ctx, query, name).Scan(
		&vmr.ID,
		&vmr.Name,
		&vmr.Namespace,
		&vmr.CPU,
		&vmr.RAM,
		&vmr.MacAddress,
		&vmr.NetworkID,
		&vmr.Image,
		&vmr.DiskSize,
		&vmr.DiskPath,
		&vmr.CreatedAt,
		&labelText,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("VM with name %q not found: %w", name, err)
		}
		return fmt.Errorf("failed to fetch VM record: %w", err)
	}

	if err := json.Unmarshal([]byte(labelText), &vmr.Labels); err != nil {
		return fmt.Errorf("failed to parse labels JSON: %w", err)
	}

	return nil
}

func (vmr *VirtualMachineRecord) GetRecordByNamespace(
	ctx context.Context,
	db *sql.DB,
	name string,
	namespace string,
) error {
	fmt.Println("-> GetRecordByNamespace")
	query := fmt.Sprintf(`
	SELECT id, name, namespace,
	       cpu, ram, mac_address,
	       network_id, image,
	       disk_size, disk_path,
	       created_at, labels
	FROM %s
	WHERE namespace = ? AND name = ? `,
		VMsTable,
	)

	var labelText string

	err := db.QueryRowContext(ctx, query, namespace, name).Scan(
		&vmr.ID,
		&vmr.Name,
		&vmr.Namespace,
		&vmr.CPU,
		&vmr.RAM,
		&vmr.MacAddress,
		&vmr.NetworkID,
		&vmr.Image,
		&vmr.DiskSize,
		&vmr.DiskPath,
		&vmr.CreatedAt,
		&labelText,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("VM %q in namespace %q not found: %w", name, namespace, err)
		}
		return fmt.Errorf("failed to retrieve VM %q: %w", name, err)
	}

	if err := json.Unmarshal([]byte(labelText), &vmr.Labels); err != nil {
		return fmt.Errorf("failed to parse VM labels: %w", err)
	}

	return nil
}
