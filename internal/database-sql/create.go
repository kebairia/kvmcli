package databasesql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// InsertVM inserts a VirtualMachineRecord into the vms table.
// It ensures that the table exists, marshals the Labels field to JSON,
// and then executes the insert statement.
func InsertVM(ctx context.Context, db *sql.DB, record *VirtualMachineRecord) error {
	if db == nil {
		return fmt.Errorf("DB is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	// Ensure the vms table exists.
	if err := EnsureVMTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure vms table exists: %w", err)
	}

	if err := EnsureNetworkTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure networks table exists: %w", err)
	}

	// Marshal the Labels map into JSON for storage in the TEXT column.
	labelsJSON, err := json.Marshal(record.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	// Define the INSERT query.
	const query = `
		INSERT INTO vms (
			name,
			namespace,
			cpu,
			ram,
			mac_address,
			network_id,
			image,
			disk_size,
			disk_path,
			created_at,
			labels
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the query using record values.
	if _, err := db.Exec(query,
		record.Name,
		record.Namespace,
		record.CPU,
		record.RAM,
		record.MacAddress,
		record.NetworkID,
		record.Image,
		record.DiskSize,
		record.DiskPath,
		record.CreatedAt,
		string(labelsJSON),
	); err != nil {
		return fmt.Errorf("failed to insert VM record: %w", err)
	}

	return nil
}
