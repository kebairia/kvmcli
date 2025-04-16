package databasesql

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

// EnsureVMTable creates the vms table if it doesn't exist.
func EnsureNetworkTable(ctx context.Context, db *sql.DB) error {
	const schema = `
  CREATE TABLE IF NOT EXISTS networks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    namespace TEXT,
    labels TEXT,
    mac_address TEXT,
    bridge TEXT,
    mode TEXT,
    net_address TEXT,
    netmask TEXT,
    dhcp TEXT,
    autostart BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

	CREATE UNIQUE INDEX IF NOT EXISTS idx_net_name_namespace ON networks(name, namespace);

	`
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create networks table: %w", err)
	}
	return nil
}
