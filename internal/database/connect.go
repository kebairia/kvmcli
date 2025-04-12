package database

import (
	"context"
	"log"
	"time"

	"github.com/kebairia/kvmcli/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx                 = context.Background()
	client              *mongo.Client
	Database            = "kvmcli"
	StoreCollection     = "store"
	VMsCollection       = "vms"
	NetworksCollection  = "networks"
	SnapshotsCollection = "snapshots"
)

func init() {
	// Build the MongoDB client options with the connection URI.
	uri := options.Client().ApplyURI("mongodb://root:example@localhost:27017")

	// Create a context with a timeout for connecting to MongoDB.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	var err error

	client, err = mongo.Connect(ctx, uri)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping the database to verify the connection.
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	logger.Log.Debugf("Connected to MongoDB")

	// Get a handle for the "kvmcli" database.
	db := client.Database("kvmcli")
	// List of collection names on which to create indexes.
	collections := []string{VMsCollection, NetworksCollection, SnapshotsCollection}
	// Define the compound index model for collections that require a unique combination of name and namespace.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", 1},
			{"namespace", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	// Loop over each collection and create the index.
	for _, collName := range collections {
		collection := db.Collection(collName)
		_, err = collection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			logger.Log.Errorf("failed to create index on collection %q: %v\n", collName, err)
		} else {
			logger.Log.Debugf("Created index on collection %q", collName)
		}
	}
	// Create an index of store object
	storeCollection := db.Collection(StoreCollection)
	storeIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = storeCollection.Indexes().CreateOne(ctx, storeIndexModel)

	if err != nil {
		logger.Log.Errorf("failed to create index on collection %v: %v\n", storeCollection, err)
	} else {
		logger.Log.Debugf("Created index on collection %v", storeCollection)
	}
}
