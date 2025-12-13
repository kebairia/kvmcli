package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Store struct {
	ID            int
	Name          string
	Namespace     string
	Labels        map[string]string
	Backend       string
	ArtifactsPath string
	ImagesPath    string
	Images        []Image
	CreatedAt     time.Time
}

type Image struct {
	ID        int64
	StoreID   int64
	Name      string
	Display   string
	Version   string
	OsProfile string
	File      string
	Checksum  string
	Size      string
	CreatedAt time.Time
}

type VMImageInfo struct {
	StoreID        int
	StoreName      string
	StoreNamespace string
	ImageID        int64
	ImageStoreID   int64
	ImageName      string
	ImageDisplay   string
	ImageVersion   string
	OsProfile      string
	ArtifactsPath  string
	ImagesPath     string
	ImageFile      string
	Checksum       string
	Size           string
}

// NewStore creates a new store record from the provided store configuration.
// func NewStore(s *store.Store) *Store {
// 	images := make(map[string]Image, len(s.Spec.Images))
// 	maps.Copy(images, s.Spec.Images)
//
// 	return &Store{
// 		Name:          s.Metadata.Name,
// 		Namespace:     s.Metadata.Namespace,
// 		Labels:        s.Metadata.Labels,
// 		Backend:       s.Spec.Backend,
// 		ArtifactsPath: s.Spec.Config.ArtifactsPath,
// 		ImagesPath:    s.Spec.Config.ImagesPath,
// 		Images:        images,
// 		CreatedAt:     time.Now(),
// 	}
// }

// EnsureVMTable creates the vms table if it doesn't exist.
func EnsureStoreTable(ctx context.Context, db *sql.DB) error {
	const schema = `
		CREATE TABLE IF NOT EXISTS ` + storesTable + ` (
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
		  ON ` + storesTable + `(name, namespace);
		
		CREATE TABLE IF NOT EXISTS ` + imagesTable + ` (
		  id         INTEGER PRIMARY KEY AUTOINCREMENT,
		  store_id   INTEGER NOT NULL,
		  name 			 TEXT,
		  display    TEXT,
		  version    TEXT,
		  os_profile TEXT,
		  file       TEXT,
		  checksum   TEXT,
		  size       TEXT,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  FOREIGN KEY(store_id) REFERENCES ` + storesTable + `(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_images_store_id_name ON ` + imagesTable + `(store_id, name);
		`
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create store table: %w", err)
	}
	return nil
}

// func getImage(
func (store *Store) GetRecord(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	const query = `
		SELECT
      id, name, namespace, backend,
      artifacts_path, images_path, created_at
		FROM ` + storesTable + `
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

// NOTE: I need to add checker for the namespace if exist no before doing the query
func (store *Store) GetRecordByNamespace(
	ctx context.Context,
	db *sql.DB,
	name, namespace string,
) error {
	const query = `
		SELECT
      id, name, namespace, backend,
      artifacts_path, images_path, created_at
		FROM ` + storesTable + `
    WHERE namespace = ? AND name = ?;
		`
	row := db.QueryRowContext(ctx, query, namespace, name)
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

func GetImage(ctx context.Context, db *sql.DB, imgName string) (*VMImageInfo, error) {
	rec := &VMImageInfo{}
	const query = `
		SELECT 
		store.id, store.name, store.namespace,
		image.id, image.store_id, image.name, image.display, image.version, image.os_profile,
		store.artifacts_path, store.images_path, image.file, image.checksum, image.size
		FROM ` + imagesTable + ` AS image
		JOIN ` + storesTable + ` AS store ON image.store_id = store.id
		WHERE image.name = ?;
		`

	row := db.QueryRowContext(ctx, query, imgName)
	if err := row.Scan(
		&rec.StoreID, &rec.StoreName, &rec.StoreNamespace,
		&rec.ImageID, &rec.ImageStoreID, &rec.ImageName, &rec.ImageDisplay, &rec.ImageVersion,
		&rec.OsProfile, &rec.ArtifactsPath, &rec.ImagesPath,
		&rec.ImageFile, &rec.Checksum, &rec.Size,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no image %q in store", imgName)
		}

		return nil, fmt.Errorf("scan image: %w", err)
	}

	return rec, nil
}

// GetStoreIDByName retrieves the database ID of a store by its name.
// Returns sql.ErrNoRows if no store with that name exists.
func GetStoreIDByName(ctx context.Context, db *sql.DB, name string) (int, error) {
	const query = `
        SELECT id FROM ` + storesTable + `
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

func (store *Store) Insert(ctx context.Context, db *sql.DB) error {
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
		INSERT INTO ` + storesTable + ` (
			name, namespace, labels,
			backend, artifacts_path, images_path
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := tx.ExecContext(ctx, storeInsert,
		store.Name,
		store.Namespace,
		string(labelsJSON),
		store.Backend,
		store.ArtifactsPath,
		store.ImagesPath,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("store %q already exists in namespace %q", store.Name, store.Namespace)
		}
		return fmt.Errorf("insert store: %w", err)
	}

	// 5. Get the new store’s primary key
	storeID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	// 6. Insert each image pointing back to the store
	const imgInsert = `
		INSERT INTO ` + imagesTable + ` (
			store_id, name, display, version, os_profile,
			file, checksum, size
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	for _, img := range store.Images {
		_, err = tx.ExecContext(ctx, imgInsert,
			storeID,
			img.Name,
			img.Display,
			img.Version,
			img.OsProfile,
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

func (store *Store) Delete(ctx context.Context, db *sql.DB) error {
	// Create a filter matching the record with the specified name
	const query = `
		DELETE FROM ` + storesTable + `
		WHERE name = ? and namespace = ?
		`

	if _, err := db.ExecContext(ctx, query, store.Name, store.Namespace); err != nil {
		return fmt.Errorf(
			"failed to delete from %s where name = %v: %w",
			storesTable,
			store.Name,
			err,
		)
	}
	return nil
}
