package database

import (
	"context"
	"database/sql"
	"fmt"
)

func Delete(ctx context.Context, db *sql.DB, name string, table string) error {
	// Create a filter matching the record with the specified name
	query := fmt.Sprintf("DELETE FROM %s WHERE name = ?", table)

	if _, err := db.ExecContext(ctx, query, name); err != nil {
		return fmt.Errorf("failed to delete from %s where name = %v: %w", table, name, err)
	}
	return nil
}
