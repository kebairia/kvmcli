package databasesql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Constants
const (
	DBFilePath     = "/home/zakaria/dox/homelab/kvmcli/kvmcli.db"
	DatabaseName   = "kvmcli"
	StoreTable     = "store"
	VMsTable       = "vms"
	NetworksTable  = "networks"
	SnapshotsTable = "snapshots"
)

// Global shared variables (only if you decide to use globals).
var (
	DB  *sql.DB
	Ctx = context.Background()
)

// InitDB opens a database handle and verifies the connection using context.
// It returns a ready-to-use *sql.DB or an error if the connection fails.
func InitDB() (*sql.DB, error) {
	// Define a 5-second timeout context for DB operations
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// Open a handle to the SQLite database (does not connect yet)
	var err error
	db, err := sql.Open("sqlite3", DBFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB handle: %w", err)
	}

	// Verify the connection with PingContext (this actually connects)
	if err := db.PingContext(Ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return db, nil
}
