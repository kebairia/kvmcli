// Package resources defines core interfaces for KVMCLI resource management.
package resources

import (
	"context"
	"database/sql"
	"text/tabwriter"

	"github.com/digitalocean/go-libvirt"
)

// Resource encapsulates operations to provision and tear down a generic KVM resource.
type Resource interface {
	// Create provisions the resource in libvirt.
	Create() error
	// Delete removes the resource from libvirt.
	Delete() error
}

// Record encapsulates database persistence operations for a resource.
type Record interface {
	// Insert persists the record into the given SQL database.
	Insert(ctx context.Context, db *sql.DB) error
	// Delete removes the record from the database.
	Delete(ctx context.Context, db *sql.DB) error
	// Get retrieves a single record by its name.
	GetRecord(ctx context.Context, db *sql.DB, name string) error
	// GetByNamespace retrieves a record by both namespace and name.
	GetRecordByNamespace(ctx context.Context, db *sql.DB, name, namespace string) error
	// Get retrieves all record for a specific table .
	// GetRecords(ctx context.Context, db *sql.DB, table string) error
	// // Get retrieves all record for a specific table .
	// GetRecordsByNamespace(ctx context.Context, db *sql.DB, namespace, table string) error
}

// ResourceInfo defines methods to render a resource’s information in a tabular, CLI-friendly format.
type ResourceInfo interface {
	// Header returns a tabwriter.Writer with column headers pre-written.
	Header() *tabwriter.Writer
	// PrintInfo writes the resource’s data as a row to the provided tabwriter.Writer.
	PrintRow(w *tabwriter.Writer)
}

// ClientSetter is implemented by any type that needs to receive a libvirt connection.
type ClientSetter interface {
	// SetConnection assigns a libvirt connection for subsequent operations.
	SetConnection(ctx context.Context, database *sql.DB, conn *libvirt.Libvirt)
}
