package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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
	Network    string
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
	CREATE INDEX IF NOT EXISTS idx_vm_name ON vms(name);
	CREATE INDEX IF NOT EXISTS idx_vm_namespace ON vms(namespace);
	`
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create vms table: %w", err)
	}
	return nil
}

// EnsureVMTable creates the vms table if it doesn't exist.
func EnsureNetworkTable(ctx context.Context, db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS networks (
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
	CREATE INDEX IF NOT EXISTS idx_vm_name ON vms(name);
	CREATE INDEX IF NOT EXISTS idx_vm_namespace ON vms(namespace);
	`
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create vms table: %w", err)
	}
	return nil
}

func InsertVM(ctx context.Context, db *sql.DB, record VirtualMachineRecord) error {
	EnsureVMTable(ctx, db)
	query := `
	INSERT INTO vms (
	name, namespace, cpu, ram, mac_address, network_id,
	image, disk_size, disk_path, created_at, labels)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.ExecContext(ctx, query,
		record.Name, record.Namespace, record.CPU, record.RAM, record.MacAddress, record.NetworkID,
		record.Image, record.DiskSize, record.DiskPath, record.CreatedAt, record.Labels,
	)
	if err != nil {
		return fmt.Errorf("failed to insert VM record: %w", err)
	}

	return nil
}
