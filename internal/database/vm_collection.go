package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VMRecord struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name"`
	Namespace   string               `bson:"namespace"`
	CPU         int                  `bson:"cpu"`
	RAM         int                  `bson:"ram"`
	MacAddress  string               `bson:"macaddress"`
	Network     string               `bson:"network"`
	SnapshotIDs []primitive.ObjectID `bson:"snapshotIds,omitempty"`
	CreatedAt   time.Time            `bson:"timestamp"`
	Disk        Disk                 `bson:"disk"`
	Image       string               `bson:"image"`
	Labels      map[string]string    `bson:"labels"`
}
type Disk struct {
	Size string `bson:"size"`
	Path string `bson:"path"`
}
type Snapshot struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	VMID    primitive.ObjectID `bson:"vmId"`
	State   string             `bson:"state"`
	TakenAt time.Time          `bson:"timestamp"`
}
type Network struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	CIDR      string             `bson:"cidr"` // example network info
	CreatedAt time.Time          `bson:"timestamp"`
}
