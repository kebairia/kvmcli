package databasesql

import (
	"context"
	"database/sql"
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
