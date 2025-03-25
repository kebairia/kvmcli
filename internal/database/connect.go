package database

import (
	"context"
	"fmt"
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

	// Create a unique compound index on "name" and "namespace"
	db := client.Database("kvmcli")
	collection := db.Collection("vms")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", 1},
			{"namespace", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Printf("failed to create index: %v", err)
	}
}
