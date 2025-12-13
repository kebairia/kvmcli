package database

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetStoreByName(t *testing.T) {
	// Setup in-memory DB
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// 1. Create tables
	if err := EnsureStoreTable(ctx, db); err != nil {
		t.Fatalf("EnsureStoreTable failed: %v", err)
	}

	// 2. Prepare test data
	testStore := Store{
		Name:          "test-store",
		Namespace:     "default",
		Labels:        map[string]string{"env": "test"},
		Backend:       "dir",
		ArtifactsPath: "/tmp/artifacts",
		ImagesPath:    "/tmp/images",
		Images: []Image{
			{
				Name:      "ubuntu",
				Version:   "20.04",
				OsProfile: "ubuntu",
				File:      "ubuntu.qcow2",
				Checksum:  "abc12345",
				Size:      "1GB",
			},
		},
		// CreatedAt will be set by DB default, so we won't strictly equate it for input
	}

	// 3. Insert store
	if err := testStore.Insert(ctx, db); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// 4. Test GetStoreByName
	retrieved, err := GetStoreByName(ctx, db, "test-store")
	if err != nil {
		t.Fatalf("GetStoreByName failed: %v", err)
	}

	// Verify fields
	if retrieved.Name != testStore.Name {
		t.Errorf("expected name %q, got %q", testStore.Name, retrieved.Name)
	}
	if retrieved.Namespace != testStore.Namespace {
		t.Errorf("expected namespace %q, got %q", testStore.Namespace, retrieved.Namespace)
	}
	if retrieved.Backend != testStore.Backend {
		t.Errorf("expected backend %q, got %q", testStore.Backend, retrieved.Backend)
	}
	if retrieved.ArtifactsPath != testStore.ArtifactsPath {
		t.Errorf("expected artifacts_path %q, got %q", testStore.ArtifactsPath, retrieved.ArtifactsPath)
	}
	if retrieved.ImagesPath != testStore.ImagesPath {
		t.Errorf("expected images_path %q, got %q", testStore.ImagesPath, retrieved.ImagesPath)
	}

	// Check labels
	if !reflect.DeepEqual(retrieved.Labels, testStore.Labels) {
		t.Errorf("expected labels %v, got %v", testStore.Labels, retrieved.Labels)
	}

	// Check Images
	if len(retrieved.Images) != len(testStore.Images) {
		t.Fatalf("expected %d images, got %d", len(testStore.Images), len(retrieved.Images))
	}
	img := retrieved.Images[0]
	expectedImg := testStore.Images[0]

	if img.Name != expectedImg.Name {
		t.Errorf("image name mismatch")
	}
	if img.Version != expectedImg.Version {
		t.Errorf("image version mismatch")
	}
	// Check store_id is set correctly (1)
	if img.StoreID == 0 {
		t.Errorf("image store_id should be set")
	}

	// Test Not Found
	_, err = GetStoreByName(ctx, db, "non-existent")
	if err == nil {
		t.Error("expected error for non-existent store, got nil")
	} else if err.Error() != `no store found with name "non-existent"` {
		t.Errorf("unexpected error message: %v", err)
	}
}
