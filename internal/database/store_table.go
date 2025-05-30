package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// keep the detailed image description
type StoreImage struct {
	Version   string `json:"version"`
	Name      string `json:"name"`
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
		CREATE TABLE IF NOT EXISTS stores (
		  id              INTEGER PRIMARY KEY AUTOINCREMENT,
		  name            TEXT    NOT NULL,
		  namespace       TEXT,
		  labels          TEXT,
		  backend         TEXT    NOT NULL,
		  artifacts_path  TEXT    NOT NULL,
		  images_path     TEXT    NOT NULL,
		  created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_stores_name_namespace
		  ON stores(name, namespace);
		
		CREATE TABLE IF NOT EXISTS images (
		  id         INTEGER PRIMARY KEY AUTOINCREMENT,
		  store_id   INTEGER NOT NULL,
			name 			 TEXT,
		  version    TEXT,
		  os_profile TEXT,
		  directory  TEXT,
		  file       TEXT,
		  checksum   TEXT,
		  size       TEXT,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY(store_id) REFERENCES stores(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_images_store_id_name ON images(store_id, name);
		`
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create store table: %w", err)
	}
	return nil
}

type ImageRecord struct {
	ID        int64
	StoreID   int64
	Name      string
	Version   string
	OsProfile string
	Directory string
	File      string
	Checksum  string
	Size      string
	CreatedAt time.Time
}

// func getImageRecord(
func (store *StoreRecord) GetRecord(
	ctx context.Context,
	db *sql.DB,
	storeID int64,
	name string,
) (*ImageRecord, error) {
	// const query = `
	// 	SELECT
	// 	id, store_id, name, version, os_profile,
	// 	directory, file, checksum, size, created_at
	// 	FROM images
	// 	WHERE store_id = ? AND name = ? ;
	// `
	const query = `
		SELECT
      i.id, i.store_id, i.name, i.version, i.os_profile,
      i.directory, i.file, i.checksum, i.size, i.created_at,
      s.id, s.name, s.namespace, s.backend,
      s.artifacts_path, s.images_path, s.created_at
    FROM images AS i
    JOIN stores AS s ON i.store_id = s.id
    WHERE i.store_id = ? AND i.name = ?;
		`

	rec := &ImageRecord{}
	row := db.QueryRowContext(ctx, query, storeID, name)
	if err := row.Scan(
		&rec.ID,
		&rec.StoreID,
		&rec.Name,
		&rec.Version,
		&rec.OsProfile,
		&rec.Directory,
		&rec.File,
		&rec.Checksum,
		&rec.Size,
		&rec.CreatedAt,
		&store.ID,
		&store.Name,
		&store.Namespace,
		&store.Backend,
		&store.ArtifactsPath,
		&store.ImagesPath,
		&store.Created_at,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no image %q in store %d", name, storeID)
		}

		return nil, fmt.Errorf("scan image: %w", err)
	}

	return rec, nil
}
