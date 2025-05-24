package database

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

func InsertNet(ctx context.Context, db *sql.DB, record *VirtualNetworkRecord) error {
	if db == nil {
		return fmt.Errorf("DB is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	// Ensure the vms table exists.
	if err := EnsureNetworkTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure %q table exists: %w", NetworksTable, err)
	}

	// Marshal the Labels map into JSON for storage in the TEXT column.
	labelsJSON, err := json.Marshal(record.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	DHCPJSON, err := json.Marshal(record.DHCP)
	if err != nil {
		return fmt.Errorf("failed to marshal DHCP: %w", err)
	}

	// Define the INSERT query.
	const query = `
		INSERT INTO networks (
			name,
			namespace,
			labels,
			mac_address,
			bridge,
			mode,
		  net_address,
		  netmask,
		  dhcp,
			autostart,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the query using record values.
	if _, err := db.Exec(query,
		record.Name,
		record.Namespace,
		string(labelsJSON),
		record.MacAddress,
		record.Bridge,
		record.Mode,
		record.NetAddress,
		record.Netmask,
		string(DHCPJSON),
		record.Autostart,
		record.CreatedAt,
	); err != nil {
		return fmt.Errorf("failed to insert Network record: %w", err)
	}

	return nil
}

func InsertStore(ctx context.Context, db *sql.DB, record *StoreRecord) error {
	if db == nil {
		return fmt.Errorf("DB is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	labelsJSON, err := json.Marshal(record.Labels)
	if err != nil {
		return err
	}
	imagesJSON, err := json.Marshal(record.Images)
	if err != nil {
		return err
	}
	// NOTE: Add EnstureStoreTable function

	// Ensure the vms table exists.
	if err := EnsureStoreTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure %q table exists: %w", NetworksTable, err)
	}
	const query = `
		INSERT INTO store (
		name, namespace, labels, backend, artifacts_path, images_path, images, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
	if _, err := db.Exec(query,
		record.Name,
		record.Namespace,
		string(labelsJSON),
		record.Backend,
		record.ArtifactsPath,
		record.ImagesPath,
		string(imagesJSON),
		record.Created_at,
	); err != nil {
		panic(err)
	}
	return nil
}
