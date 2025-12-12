package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

func (vm *VirtualMachine) ScanRows(rows *sql.Rows) error {
	var labelJSON string
	if err := rows.Scan(
		&vm.ID, &vm.Name, &vm.Namespace,
		&vm.CPU, &vm.RAM, &vm.MacAddress,
		&vm.NetworkID, &vm.Image, &vm.DiskSize,
		&vm.DiskPath, &vm.CreatedAt, &labelJSON,
	); err != nil {
		return err
	}

	return json.Unmarshal([]byte(labelJSON), &vm.Labels)
}

func (vm *VirtualMachine) ScanRow(row *sql.Row) error {
	var labelJSON string
	if err := row.Scan(
		&vm.ID, &vm.Name, &vm.Namespace,
		&vm.CPU, &vm.RAM, &vm.MacAddress,
		&vm.NetworkID, &vm.Image, &vm.DiskSize,
		&vm.DiskPath, &vm.CreatedAt, &labelJSON,
	); err != nil {
		return err
	}

	return json.Unmarshal([]byte(labelJSON), &vm.Labels)
}

// Network records

func (net *VirtualNetwork) ScanRows(rows *sql.Rows) error {
	var labelJSON string
	var DHCPJSON string
	if err := rows.Scan(
		&net.ID, &net.Name, &net.Namespace,
		&labelJSON, &net.MacAddress, &net.Bridge,
		&net.Mode, &net.NetAddress, &net.Netmask,
		&DHCPJSON, &net.Autostart, &net.CreatedAt,
	); err != nil {
		return err
	}

	// Decode labels
	if err := json.Unmarshal([]byte(labelJSON), &net.Labels); err != nil {
		return fmt.Errorf("failed to unmarshal labels: %w", err)
	}

	// Decode DHCP
	// if err := json.Unmarshal([]byte(DHCPJSON), &net.DHCP); err != nil {
	// 	return fmt.Errorf("failed to unmarshal DHCP: %w", err)
	// }

	return nil
}

func (net *VirtualNetwork) ScanRow(row *sql.Row) error {
	var labelJSON string
	var DHCPJSON string
	if err := row.Scan(
		&net.ID, &net.Name, &net.Namespace,
		&labelJSON, &net.MacAddress, &net.Bridge,
		&net.Mode, &net.NetAddress, &net.Netmask,
		&DHCPJSON, &net.Autostart, &net.CreatedAt,
	); err != nil {
		return err
	}

	// Decode labels
	if err := json.Unmarshal([]byte(labelJSON), &net.Labels); err != nil {
		return fmt.Errorf("failed to unmarshal labels: %w", err)
	}

	// Decode DHCP
	// if err := json.Unmarshal([]byte(DHCPJSON), &net.DHCP); err != nil {
	// 	return fmt.Errorf("failed to unmarshal DHCP: %w", err)
	// }

	return nil
}

// other
func GetNetworkNameByID(ctx context.Context, db *sql.DB, id int) (string, error) {
	const query = `SELECT name FROM networks WHERE id = ?`
	var name string

	err := db.QueryRowContext(ctx, query, id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no network found with ID %d", id)
		}
		return "", fmt.Errorf("faild to fetch network name for ID %d: %w", id, err)
	}
	return name, nil
}

func GetNetworkIDByName(ctx context.Context, db *sql.DB, networkName string) (int, error) {
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
