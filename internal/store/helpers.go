package store

import (
	"maps"
	"time"

	db "github.com/kebairia/kvmcli/internal/database"
)

// NOTE: I used this before, I need to check it later
//
//	for name, img := range s.Spec.Images {
//		images[name] = img
//	}
//
// NewStoreRecord creates a new store record from the provided store configuration.
func NewStoreRecord(s *Store) *db.StoreRecord {
	images := make(map[string]db.ImageRecord, len(s.Spec.Images))
	maps.Copy(images, s.Spec.Images)

	return &db.StoreRecord{
		Name:          s.Metadata.Name,
		Namespace:     s.Metadata.Namespace,
		Labels:        s.Metadata.Labels,
		Backend:       s.Spec.Backend,
		ArtifactsPath: s.Spec.Config.ArtifactsPath,
		ImagesPath:    s.Spec.Config.ImagesPath,
		Images:        images,
		CreatedAt:     time.Now(),
	}
}
