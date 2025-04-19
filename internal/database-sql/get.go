package databasesql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// GetObjectsByNamespace retrieves all documents of type T from the specified collection
// that match the given namespace.
func BringObjectsByNamespace(
	ctx context.Context,
	db *sql.DB,
	namespace, table string,
) ([]VirtualMachineRecord, error) {
	query := fmt.Sprintf(`
		SELECT id, name, namespace, 
		       cpu, ram, mac_address, 
		       network_id, image, 
		       disk_size, disk_path, 
		       created_at, labels
		FROM %s
		WHERE namespace = ?`, table)
	var (
		objects   []VirtualMachineRecord
		labelText string
	)
	rows, err := db.QueryContext(ctx, query, namespace)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var object VirtualMachineRecord
		if err := rows.Scan(
			&object.ID,
			&object.Name,
			&object.Namespace,
			&object.CPU,
			&object.RAM,
			&object.MacAddress,
			&object.NetworkID,
			&object.Image,
			&object.DiskSize,
			&object.DiskPath,
			&object.CreatedAt,
			&labelText,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		objects = append(objects, object)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return objects, nil
}

func BringRecord(
	ctx context.Context,
	db *sql.DB,
	name, table string,
) (*VirtualMachineRecord, error) {
	// Only allow known tables (basic safety)
	if table != "vms" {
		return nil, fmt.Errorf("table %q is not allowed", table)
	}

	query := fmt.Sprintf(`
		SELECT id, name, namespace, 
		       cpu, ram, mac_address, 
		       network_id, image, 
		       disk_size, disk_path, 
		       created_at, labels
		FROM %s
		WHERE name = ?`, table)

	var (
		record    VirtualMachineRecord
		labelText string
	)

	err := db.QueryRowContext(ctx, query, name).Scan(
		&record.ID,
		&record.Name,
		&record.Namespace,
		&record.CPU,
		&record.RAM,
		&record.MacAddress,
		&record.NetworkID,
		&record.Image,
		&record.DiskSize,
		&record.DiskPath,
		&record.CreatedAt,
		&labelText,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("VM with name %q not found: %w", name, err)
		}
		return nil, fmt.Errorf("failed to fetch VM record: %w", err)
	}

	if err := json.Unmarshal([]byte(labelText), &record.Labels); err != nil {
		return nil, fmt.Errorf("failed to parse labels JSON: %w", err)
	}

	return &record, nil
}

func GetNetworkRecord(
	ctx context.Context,
	db *sql.DB,
	name, table string,
) (*VirtualNetworkRecord, error) {
	// Only allow known tables (basic safety)
	if table != "networks" {
		return nil, fmt.Errorf("table %q is not allowed", table)
	}

	query := fmt.Sprintf(`
		SELECT id, name, namespace, 
		       labels, mac_address, 
		       bridge, mode, 
		       net_address, netmask, 
		       dhcp, autostart, created_at
		FROM %s
		WHERE name = ?`, table)

	var (
		record    VirtualNetworkRecord
		labelText string
		DHCPText  string
	)

	err := db.QueryRowContext(ctx, query, name).Scan(
		&record.ID,
		&record.Name,
		&record.Namespace,
		&labelText,
		&record.MacAddress,
		&record.Bridge,
		&record.Mode,
		&record.NetAddress,
		&record.Netmask,
		&DHCPText,
		&record.Autostart,
		&record.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("VM with name %q not found: %w", name, err)
		}
		return nil, fmt.Errorf("failed to fetch VM record: %w", err)
	}

	if err := json.Unmarshal([]byte(DHCPText), &record.DHCP); err != nil {
		return nil, fmt.Errorf("failed to parse labels JSON: %w", err)
	}
	if err := json.Unmarshal([]byte(labelText), &record.Labels); err != nil {
		return nil, fmt.Errorf("failed to parse labels JSON: %w", err)
	}

	return &record, nil
}

func GetNetworkObjectsByNamespace(
	ctx context.Context,
	db *sql.DB,
	namespace, table string,
) ([]VirtualNetworkRecord, error) {
	query := fmt.Sprintf(`
		SELECT id, name, namespace, 
		       labels, mac_address, 
		       bridge, mode, 
		       net_address, netmask, 
		       dhcp, autostart, created_at
		FROM %s
		WHERE namespace = ?`, table)

	var (
		objects   []VirtualNetworkRecord
		labelText string
		DHCPText  string
	)
	rows, err := db.QueryContext(ctx, query, namespace)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var object VirtualNetworkRecord
		if err := rows.Scan(
			&object.ID,
			&object.Name,
			&object.Namespace,
			&labelText,
			&object.MacAddress,
			&object.Bridge,
			&object.Mode,
			&object.NetAddress,
			&object.Netmask,
			&DHCPText,
			&object.Autostart,
			&object.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		objects = append(objects, object)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return objects, nil
}
