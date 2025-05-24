package store

import (
	"time"

	db "github.com/kebairia/kvmcli/internal/database"
)

// NewStoreRecord converts a Store into a database StoreRecord.
//
//	func NewStoreRecordold(s *Store) *db.StoreRecord {
//		// Convert a slice of StoreImage to a slice of database image
//		images := make(map[string]db.StoreImage, len(s.Spec.Images))
//		for dist, img := range s.Spec.Images {
//			images[dist] = db.StoreImage{
//				Version:   img.Version,
//				OsProfile: img.OsProfile,
//				Directory: img.Directory,
//				File:      img.File,
//				Checksum:  img.Checksum,
//				Size:      img.Size,
//			}
//		}
//		return &db.StoreRecord{
//			Metadata: db.StoreMetadata{
//				Name:      s.Metadata.Name,
//				Namespace: s.Metadata.Namespace,
//				Labels:    s.Metadata.Labels,
//				CreatedAt: time.Now(),
//			},
//			Spec: db.StoreSpec{
//				Backend: s.Spec.Backend,
//				Config: db.StoreConfig{
//					ArtifactsPath: s.Spec.Config.ArtifactsPath,
//					ImagesPath:    s.Spec.Config.ImagesPath,
//				},
//				Images: images,
//			},
//		}
//	}
func NewStoreRecord(s *Store) *db.StoreRecord {
	images := make(map[string]db.StoreImage, len(s.Spec.Images))

	for name, img := range s.Spec.Images {
		images[name] = img // keep it structured
	}
	return &db.StoreRecord{
		Name:          s.Metadata.Name,
		Namespace:     s.Metadata.Namespace,
		Labels:        s.Metadata.Labels,
		Backend:       s.Spec.Backend,
		ArtifactsPath: s.Spec.Config.ArtifactsPath,
		ImagesPath:    s.Spec.Config.ImagesPath,
		Images:        images,
		Created_at:    time.Now(),
	}
}
