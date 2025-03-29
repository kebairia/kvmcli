package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NetRecord struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Namespace  string             `bson:"namespace"`
	Labels     map[string]string  `bson:"labels"`
	MacAddress string             `bson:"macaddress"`
	Bridge     string             `bson:"bridge"`
	Mode       string             `bson:"mode"`
	NetAddress string             `bson:"netAddress"`
	Netmask    string             `bson:"netmask"`
	DHCP       map[string]string  `bson:"dhcp"`
	CreatedAt  time.Time          `bson:"timestamp"`
}
