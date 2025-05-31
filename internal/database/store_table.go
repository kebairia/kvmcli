package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type StoreRecord struct {
	ID            int
	Name          string
	Namespace     string
	Labels        map[string]string
	Backend       string
	ArtifactsPath string
	ImagesPath    string
	Images        map[string]ImageRecord
	CreatedAt     time.Time
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

// func getImageRecord(
func (store *StoreRecord) GetRecord(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	const query = `
		SELECT
      id, name, namespace, backend,
      artifacts_path, images_path, created_at
    WHERE name = ?;
		`
	row := db.QueryRowContext(ctx, query, name)
	if err := row.Scan(
		&store.ID,
		&store.Name,
		&store.Namespace,
		&store.Backend,
		&store.ArtifactsPath,
		&store.ImagesPath,
		&store.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no store record for  %q ", name)
		}
	}
	return nil
}

// func getImageRecord(
func (store *StoreRecord) GetImageRecord(
	ctx context.Context,
	db *sql.DB,
	imageName string,
) (*ImageRecord, error) {
	const query = `
		SELECT
      image.id, image.store_id, image.name, image.version, image.os_profile,
      image.directory, image.file, image.checksum, image.size, image.created_at,
      store.id, store.name, store.namespace, store.backend,
      store.artifacts_path, store.images_path, store.created_at
    FROM images AS image
    JOIN stores AS store ON image.store_id = store.id
    WHERE image.store_id = ? AND image.name = ?;
		`
	rec := &ImageRecord{}
	row := db.QueryRowContext(ctx, query, store.ID, imageName)
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
		&store.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no image %q in store %d", imageName, store.ID)
		}

		return nil, fmt.Errorf("scan image: %w", err)
	}

	return rec, nil
}

// GetStoreIDByName retrieves the database ID of a store by its name.
// Returns sql.ErrNoRows if no store with that name exists.
func GetStoreIDByName(ctx context.Context, db *sql.DB, name string) (int, error) {
	const query = `
        SELECT id
        FROM stores
        WHERE name = ?
    `

	var id int
	err := db.QueryRowContext(ctx, query, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("store %q not found", name)
		}
		return 0, fmt.Errorf("query store ID by name: %w", err)
	}
	return id, nil
}

func (store *StoreRecord) Insert(ctx context.Context, db *sql.DB) error {
	// 1. Ensure tables exist (including the images table!)
	if err := EnsureStoreTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure schema: %w", err)
	}

	// 2. Marshal labels if you still store them as JSON
	labelsJSON, err := json.Marshal(store.Labels)
	if err != nil {
		return fmt.Errorf("marshal labels: %w", err)
	}

	// 3. Start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	// In case of any error, roll back
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 4. Insert the store row (no more images JSON here)
	const storeInsert = `
		INSERT INTO stores (
			name, namespace, labels,
			backend, artifacts_path, images_path
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := tx.ExecContext(ctx, storeInsert,
		store.Name,
		// record.Namespace,
		"k8s",
		string(labelsJSON),
		store.Backend,
		store.ArtifactsPath,
		store.ImagesPath,
	)
	if err != nil {
		return fmt.Errorf("insert store: %w", err)
	}

	// 5. Get the new store’s primary key
	storeID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	// 6. Insert each image pointing back to the store
	const imgInsert = `
		INSERT INTO images (
			store_id, name, version, os_profile,
			directory, file, checksum, size
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	for _, img := range store.Images {
		_, err = tx.ExecContext(ctx, imgInsert,
			storeID,
			img.Name,
			img.Version,
			"http://rockylinux.org/rocky/9",
			img.Directory,
			img.File,
			img.Checksum,
			img.Size,
		)
		if err != nil {
			return fmt.Errorf("insert image %v: %w", img, err)
		}
	}

	// 7. Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit: %w", err)
	}
	return nil
}
