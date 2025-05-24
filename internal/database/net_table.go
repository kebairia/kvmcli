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
	const query = `
		SELECT id, name, namespace,
		labels, mac_address, 
		bridge, mode, 
		net_address, netmask,
		dhcp, autostart, created_at
		FROM networks WHERE name = ?
		`
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
