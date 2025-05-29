package resources

import (
	"context"
	"database/sql"
	"text/tabwriter"

	"github.com/digitalocean/go-libvirt"
)

// Resource defines operations for managing a resource.
type Resource interface {
	Create() error
	Delete() error
	// Get(name string) (Record, error)
	// Header() *tabwriter.Writer
}

// Record defines operations for managing a records on database.
type Record interface {
	Insert(ctx context.Context, db *sql.DB) error
	Delete() error
	GetRecord(ctx context.Context, db *sql.DB, name string) error
	// ScanRow(row *sql.Row) error
	// ScanRows(rows *sql.Rows) error
}

type ResourceInfo interface {
	Header() *tabwriter.Writer
	PrintInfo(w *tabwriter.Writer)
}

// ClientSetter is implemented by types that require a libvirt connection.
type ClientSetter interface {
	SetConnection(*libvirt.Libvirt)
}
