package store

import (
	"time"

	"github.com/kebairia/kvmcli/internal/database"
	db "github.com/kebairia/kvmcli/internal/database"
)

// NewStoreRecord converts a Store into a database StoreRecord.
func NewStoreRecord(s *Store) *db.StoreRecord {
	// Convert a slice of StoreImage to a slice of database image
	images := make(map[string]database.StoreImage, len(s.Spec.Images))
	for dist, img := range s.Spec.Images {
		images[dist] = database.StoreImage{
			Version:   img.Version,
			Directory: img.Directory,
			File:      img.File,
			Checksum:  img.Checksum,
			Size:      img.Size,
		}
	}
	return &db.StoreRecord{
		Metadata: db.StoreMetadata{
			Name:      s.Metadata.Name,
			Namespace: s.Metadata.Namespace,
			Labels:    s.Metadata.Labels,
			CreatedAt: time.Now(),
		},
		Spec: db.StoreSpec{
			Backend: s.Spec.Backend,
			Config: db.StoreConfig{
				Path: s.Spec.Config.Path,
			},
			Images: images,
		},
	}
}
