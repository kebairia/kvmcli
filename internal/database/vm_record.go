package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// VMRecord represents a VM stored in the SQLite database.
type VirtualMachine struct {
	ID         int
	Name       string
	Namespace  string
	Labels     map[string]string
	CPU        int
	RAM        int
	IP         string
	MacAddress string
	NetworkID  int
	StoreID    int
	Image      string
	DiskSize   string
	DiskPath   string
	CreatedAt  time.Time
	// SnapshotIDs []string we don't use snapshot id here, in the snapshot table we reference  t the vm
}

// EnsureVMTable creates the 'vms' table and its unique index if they do not already exist.
func EnsureVMTable(ctx context.Context, db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS vms (
	  id          INTEGER PRIMARY KEY AUTOINCREMENT,
	  name        TEXT NOT NULL,
	  namespace   TEXT,
	  cpu         INTEGER,
	  ram         INTEGER,
	  ip_address  TEXT,
	  mac_address TEXT,
	  network_id  INTEGER,
	  store_id    INTEGER,
	  image       TEXT,
	  disk_size   TEXT,
	  disk_path   TEXT,
	  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	  labels      TEXT,
	  FOREIGN KEY (network_id) REFERENCES networks(id),
	  FOREIGN KEY (store_id)    REFERENCES stores(id)
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_vm_name_namespace
	  ON vms(name, namespace);
	`

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("EnsureVMTable: failed to create table/index: %w", err)
	}
	return nil
}

func (vmr *VirtualMachine) GetRecord(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	query := fmt.Sprintf(`
		SELECT id, name, namespace, 
		       cpu, ram, ip_address, mac_address, 
		       network_id, store_id, image, 
		       disk_size, disk_path, 
		       created_at, labels
		FROM %s
		WHERE name = ?`, vmsTable)

	// record    VirtualMachine
	var labelText string

	err := db.QueryRowContext(ctx, query, name).Scan(
		&vmr.ID,
		&vmr.Name,
		&vmr.Namespace,
		&vmr.CPU,
		&vmr.RAM,
		&vmr.IP,
		&vmr.MacAddress,
		&vmr.NetworkID,
		&vmr.StoreID,
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

func (vmr *VirtualMachine) GetRecordByNamespace(
	ctx context.Context,
	db *sql.DB,
	name string,
	namespace string,
) error {
	fmt.Println("-> GetRecordByNamespace")
	query := fmt.Sprintf(`
	SELECT id, name, namespace,
	       cpu, ram, ip_address, mac_address,
	       network_id, store_id, image,
	       disk_size, disk_path,
	       created_at, labels
	FROM %s
	WHERE namespace = ? AND name = ? `,
		vmsTable,
	)

	var labelText string

	err := db.QueryRowContext(ctx, query, namespace, name).Scan(
		&vmr.ID,
		&vmr.Name,
		&vmr.Namespace,
		&vmr.CPU,
		&vmr.RAM,
		&vmr.IP,
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

// Insert inserts a VirtualMachine into the vms table.
// It ensures that the table exists, marshals the Labels field to JSON,
// and then executes the insert statement.
func (vmr *VirtualMachine) Insert(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("DB is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	// Ensure the vms/networks/stores tables exists.
	if err := EnsureVMTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure vms table exists: %w", err)
	}

	if err := EnsureNetworkTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure networks table exists: %w", err)
	}

	if err := EnsureStoreTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure stores table exists: %w", err)
	}

	// Marshal the Labels map into JSON for storage in the TEXT column.
	labelsJSON, err := json.Marshal(vmr.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}
	// Define the INSERT query.
	const query = `
		INSERT INTO ` + vmsTable + ` (
			name,
			namespace,
			cpu,
			ram,
			ip_address,
			mac_address,
			network_id,
			store_id,
			image,
			disk_size,
			disk_path,
			created_at,
			labels
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

	// Execute the query using record values.
	if _, err := db.Exec(query,
		vmr.Name,
		vmr.Namespace,
		vmr.CPU,
		vmr.RAM,
		vmr.IP,
		vmr.MacAddress,
		vmr.NetworkID,
		vmr.StoreID,
		vmr.Image,
		vmr.DiskSize,
		vmr.DiskPath,
		vmr.CreatedAt,
		string(labelsJSON),
	); err != nil {
		return fmt.Errorf("failed to insert VM record: %w", err)
	}

	return nil
}

// Delete removes a network row by name+namespace.
func (v *VirtualMachine) Delete(ctx context.Context, db *sql.DB) error {
	const stmt = `
	DELETE FROM ` + vmsTable + `
	WHERE name = ? AND namespace = ?
	`
	_, err := db.ExecContext(ctx, stmt, v.Name, v.Namespace)
	if err != nil {
		return fmt.Errorf("delete vm: %w", err)
	}
	return nil
}
