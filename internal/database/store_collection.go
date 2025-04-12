package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoreRecord struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Metadata StoreMetadata      `bson:"metadata"`
	Spec     StoreSpec          `bson:"spec"`
}

type StoreMetadata struct {
	Name      string            `bson:"name"`
	Namespace string            `bson:"namespace,omitempty"`
	Labels    map[string]string `bson:"labels"`
	CreatedAt time.Time         `bson:"timestamp"`
}

type StoreSpec struct {
	Backend string                `bson:"backend"`
	Config  StoreConfig           `bson:"config"`
	Images  map[string]StoreImage `bson:"images"` // Changed from a slice to a map
}

type StoreConfig struct {
	ArtifactsPath string `bson:"artifactsPath"`
	ImagesPath    string `bson:"imagesPath"`
}

type StoreImage struct {
	Version   string `bson:"version"`
	OsProfile string `bson:"osprofile"`
	Directory string `bson:"directory"` // New field for the directory portion
	File      string `bson:"file"`
	Checksum  string `bson:"checksum"`
	Size      string `bson:"size"`
}
