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
	ctx    = context.Background()
	client *mongo.Client
)

func init() {
	uri := options.Client().ApplyURI("mongodb://root:example@localhost:27017")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	var err error

	client, err = mongo.Connect(ctx, uri)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	logger.Log.Debugf("Connected to MongoDB")

	// Get the database handle
	db := client.Database("kvmcli")
	// List of collectoins to create the index on.
	collection := []string{"vms", "networks", " snapshots"}
	// Define the compound index model.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", 1},
			{"namespace", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	// Loop over each collection and create the index.
	for _, collName := range collection {
		collection := db.Collection(collName)
		_, err = collection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			logger.Log.Errorf("failed to create index on collection %q: %v\n", collName, err)
		} else {
			logger.Log.Debugf("Created index on collection %q", collName)
		}
	}
}
