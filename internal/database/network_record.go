package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type VirtualNetworkRecord struct {
	ID         int
	Name       string
	Namespace  string
	Labels     map[string]string
	MacAddress string
	Bridge     string
	Mode       string
	NetAddress string
	Netmask    string
	DHCP       map[string]string
	Autostart  bool
	CreatedAt  time.Time
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

func (net *VirtualNetworkRecord) GetRecord(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	query := fmt.Sprintf(`
		SELECT id, name, namespace,
		labels, mac_address, 
		bridge, mode, 
		net_address, netmask,
		dhcp, autostart, created_at
		FROM %s WHERE name = ?`,
		NetworksTable,
	)

	var (
		labelText string
		DHCPText  string
	)
	err := db.QueryRowContext(ctx, query, name).Scan(
		&net.ID,
		&net.Name,
		&net.Namespace,
		&labelText,
		&net.MacAddress,
		&net.Bridge,
		&net.Mode,
		&net.NetAddress,
		&net.Netmask,
		&DHCPText,
		&net.Autostart,
		&net.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("VM with name %q not found: %w", name, err)
		}

		return fmt.Errorf("failed to fetch VM record: %w", err)
	}
	if err := json.Unmarshal([]byte(DHCPText), &net.DHCP); err != nil {
		return fmt.Errorf("failed to parse labels JSON: %w", err)
	}
	if err := json.Unmarshal([]byte(labelText), &net.Labels); err != nil {
		return fmt.Errorf("failed to parse labels JSON: %w", err)
	}
	return nil
}

func (net *VirtualNetworkRecord) GetRecordByNamespace(
	ctx context.Context,
	db *sql.DB,
	name string,
	namespace string,
) error {
	query := fmt.Sprintf(`
		SELECT id, name, namespace,
		labels, mac_address, 
		bridge, mode, 
		net_address, netmask,
		dhcp, autostart, created_at
		FROM %s WHERE namespace = ? AND name = ?`,
		NetworksTable,
	)

	var (
		labelText string
		DHCPText  string
	)
	err := db.QueryRowContext(ctx, query, namespace, name).Scan(
		&net.ID,
		&net.Name,
		&net.Namespace,
		&labelText,
		&net.MacAddress,
		&net.Bridge,
		&net.Mode,
		&net.NetAddress,
		&net.Netmask,
		&DHCPText,
		&net.Autostart,
		&net.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf(
				"network with name %q in namespace %q not found: %w",
				name,
				namespace,
				err,
			)
		}

		return fmt.Errorf("failed to fetch VM record: %w", err)
	}
	if err := json.Unmarshal([]byte(DHCPText), &net.DHCP); err != nil {
		return fmt.Errorf("failed to parse labels JSON: %w", err)
	}
	if err := json.Unmarshal([]byte(labelText), &net.Labels); err != nil {
		return fmt.Errorf("failed to parse labels JSON: %w", err)
	}
	return nil
}

func (net *VirtualNetworkRecord) Insert(ctx context.Context, db *sql.DB) error {
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
	labelsJSON, err := json.Marshal(net.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	DHCPJSON, err := json.Marshal(net.DHCP)
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
		net.Name,
		net.Namespace,
		string(labelsJSON),
		net.MacAddress,
		net.Bridge,
		net.Mode,
		net.NetAddress,
		net.Netmask,
		string(DHCPJSON),
		net.Autostart,
		net.CreatedAt,
	); err != nil {
		return fmt.Errorf("failed to insert Network record: %w", err)
	}

	return nil
}

func (net *VirtualNetworkRecord) Delete(ctx context.Context, db *sql.DB) error {
	// Create a filter matching the record with the specified name
	query := fmt.Sprintf("DELETE FROM %s WHERE name = ?", NetworksTable)

	if _, err := db.ExecContext(ctx, query, net.Name); err != nil {
		return fmt.Errorf(
			"failed to delete from %s where name = %v: %w",
			NetworksTable,
			net.Name,
			err,
		)
	}
	return nil
}
