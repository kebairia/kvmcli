package resources

import (
	"context"
	"database/sql"

	"github.com/digitalocean/go-libvirt"
)

// Resource defines operations for managing a resource.
type Resource interface {
	Create() error
	Delete() error
	// Header() *tabwriter.Writer
}

// Record defines operations for managing a records on database.
type Record interface {
	// Insert() error
	// Delete() error
	// ScanRow(row *sql.Row) error
	// ScanRows(rows *sql.Rows) error
	GetRecord(ctx context.Context, db *sql.DB, name string)
}

// ClientSetter is implemented by types that require a libvirt connection.
type ClientSetter interface {
	SetConnection(*libvirt.Libvirt)
}
