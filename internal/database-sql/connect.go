package databasesql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Database file name
const dbFileName = "/home/zakaria/dox/homelab/kvmcli/kvmcli.db"

var (
	DB  *sql.DB
	Ctx context.Context
)

// InitDB opens a database handle and verifies the connection using context.
// It returns a ready-to-use *sql.DB or an error if the connection fails.
func InitDB() (*sql.DB, error) {
	// Define a 5-second timeout context for DB operations
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	Ctx = context.Background() // ← Persistent context for app lifetime

	// Open a handle to the SQLite database (does not connect yet)
	var err error
	DB, err = sql.Open("sqlite3", dbFileName)
	if DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open DB handle: %w", err)
	}

	// Verify the connection with PingContext (this actually connects)
	if err := DB.PingContext(Ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	fmt.Println("Connected successfully ")

	return DB, nil
}
