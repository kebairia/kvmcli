package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

const (
	vmColumns      = `id, name, namespace, cpu, ram, mac_address, network_id, image, disk_size, disk_path, created_at, labels`
	networkColumns = `id, name, namespace, mac_address, bridge, mode, net_address, netmask, dhcp, autostart, created_at, labels`
)

// GetRecords retrieves all documents of type T from the specified collection
// that match the given namespace.
func GetVMRecords(
	ctx context.Context,
	db *sql.DB,
	namespace string,
) ([]VirtualMachineRecord, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", vmColumns, VMsTable)
	args := []any{}
	if namespace != "" {
		query += " WHERE namespace = ?"
		args = append(args, namespace)
	}
	var (
		vms       []VirtualMachineRecord
		rawLabels string
	)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var vm VirtualMachineRecord
		if err := rows.Scan(
			&vm.ID,
			&vm.Name,
			&vm.Namespace,
			&vm.CPU,
			&vm.RAM,
			&vm.MacAddress,
			&vm.NetworkID,
			&vm.Image,
			&vm.DiskSize,
			&vm.DiskPath,
			&vm.CreatedAt,
			&rawLabels,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		// parse JSON labels
		if err := json.Unmarshal([]byte(rawLabels), &vm.Labels); err != nil {
			return nil, fmt.Errorf("invalid labels JSON: %w", err)
		}
		vms = append(vms, vm)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return vms, nil
}

func GetNetworkRecords(
	ctx context.Context,
	db *sql.DB,
	namespace string,
) ([]VirtualNetworkRecord, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", networkColumns, NetworksTable)
	args := []any{}
	if namespace != "" {
		query += " WHERE namespace = ?"
		args = append(args, namespace)
	}
	var (
		networks           []VirtualNetworkRecord
		rawLabels, rawDHCP string
	)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var network VirtualNetworkRecord
		if err := rows.Scan(
			&network.ID,
			&network.Name,
			&network.Namespace,
			&network.MacAddress,
			&network.Bridge,
			&network.Mode,
			&network.NetAddress,
			&network.Netmask,
			&rawDHCP,
			&network.Autostart,
			&network.CreatedAt,
			&rawLabels,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		fmt.Println(network)
		networks = append(networks, network)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return networks, nil
}
