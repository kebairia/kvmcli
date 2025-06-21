package store

import (
	"time"

	db "github.com/kebairia/kvmcli/internal/database"
)

// NOTE: I need to move this to store package
// but before that, I need to remove the dependency of using
// db.Ctx and db.DB from the creation of Insert method for store struct
// I need to do that using the operator
// otherewise, I will have an import cycle when using this constructor
// on database package

// NewStoreRecord creates a new store record from the provided store configuration.
func NewStoreRecord(s *Store) *db.StoreRecord {
	images := make([]db.ImageRecord, len(s.Config.Spec.Images))

	for index, img := range s.Config.Spec.Images {
		images[index] = db.ImageRecord{
			Name:      img.Name,
			Version:   img.Version,
			OsProfile: img.OsProfile,
			File:      img.File,
			Checksum:  img.Checksum,
			Size:      img.Size,
		}
	}

	return &db.StoreRecord{
		Name:          s.Config.Metadata.Name,
		Namespace:     s.Config.Metadata.Namespace,
		Labels:        s.Config.Metadata.Labels,
		Backend:       s.Config.Spec.Backend,
		ArtifactsPath: s.Config.Spec.Paths.ArtifactsPath,
		ImagesPath:    s.Config.Spec.Paths.ImagesPath,
		Images:        images,
		CreatedAt:     time.Now(),
	}
}
