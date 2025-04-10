package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Database file name
const dbFileName = "kvmcli.db"

// InitDB opens a database handle and verifies the connection using context.
// It returns a ready-to-use *sql.DB or an error if the connection fails.
func InitDB() (*sql.DB, error) {
	// Define a 5-second timeout context for DB operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Open a handle to the SQLite database (does not connect yet)
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB handle: %w", err)
	}

	// Verify the connection with PingContext (this actually connects)
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return db, nil
}
