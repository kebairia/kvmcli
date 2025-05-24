package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// type StoreRecord struct {
// 	ID         int
// 	Name       string
// 	Namespace  string
// 	Labels     map[string]string
// 	Backend    string
// 	Config     StoreConfig
// 	Images     map[string]StoreImage
// 	Created_at time.Time
// }
// type StoreConfig struct {
// 	ArtifactsPath string
// 	ImagesPath    string
// }
// type StoreImage struct {
// 	Version   string
// 	OsProfile string
// 	Directory string
// 	File      string
// 	Checksum  string
// 	Size      string
// }

// type StoreRecord struct {
// 	ID            int
// 	Name          string
// 	Namespace     string
// 	Labels        map[string]string
// 	Backend       string
// 	ArtifactsPath string
// 	ImagesPath    string
// 	Images        map[string]string
// 	Created_at    time.Time
// }

// keep the detailed image description
type StoreImage struct {
	Version   string `json:"version"`
	OsProfile string `json:"osProfile"`
	Directory string `json:"directory"`
	File      string `json:"file"`
	Checksum  string `json:"checksum"`
	Size      string `json:"size"`
}

type StoreRecord struct {
	ID            int
	Name          string
	Namespace     string
	Labels        map[string]string
	Backend       string
	ArtifactsPath string
	ImagesPath    string
	Images        map[string]StoreImage // <-- change here
	Created_at    time.Time
}

// EnsureVMTable creates the vms table if it doesn't exist.
func EnsureStoreTable(ctx context.Context, db *sql.DB) error {
	const schema = `
  CREATE TABLE IF NOT EXISTS store (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    namespace TEXT,
    labels TEXT,
    backend TEXT,
    artifacts_path TEXT,
    images_path TEXT,
    images TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

	CREATE UNIQUE INDEX IF NOT EXISTS idx_net_name_namespace ON store(name, namespace);

	`
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create store table: %w", err)
	}
	return nil
}

func (store *StoreRecord) GetRecord(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	const query = `
		SELECT id, name, namespace,
		labels, backend, 
		artifacts_path, images_path, 
		images, created_at
		FROM store WHERE name = ?
		`
	var labelsJSON, imagesJSON string
	row := db.QueryRowContext(ctx, query, name)
	if err := row.Scan(
		&store.ID,
		&store.Name,
		&store.Namespace,
		&labelsJSON,
		&store.Backend,
		&store.ArtifactsPath,
		&store.ImagesPath,
		&imagesJSON,
		&store.Created_at,
	); err != nil {
		return fmt.Errorf("failed to get store record: %w", err)
	}
	if err := json.Unmarshal([]byte(labelsJSON), &store.Labels); err != nil {
		return fmt.Errorf("failed to unmarshal labels: %w", err)
	}
	if err := json.Unmarshal([]byte(imagesJSON), &store.Images); err != nil {
		return fmt.Errorf("failed to unmarshal images: %w", err)
	}

	return nil
}
