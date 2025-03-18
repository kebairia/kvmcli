package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VM struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name"`
	Namespace   string               `bson:"namespace"`
	CPU         int                  `bson:"cpu"`
	RAM         int                  `bson:"ram"`
	MacAddress  string               `bson:"macaddress"`
	NetworkID   string               `bson:"network"`
	SnapshotIDs []primitive.ObjectID `bson:"snapshotIds,omitempty"`
	CreatedAt   time.Time            `bson:"timestamp"`
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
