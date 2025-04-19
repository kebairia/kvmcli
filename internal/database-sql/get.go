package databasesql

import (
	"context"
	"database/sql"
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
