package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx    = context.TODO()
	client *mongo.Client
)

func init() {
	uri := options.Client().ApplyURI("mongodb://root:example@localhost:27017")

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
	log.Println("Connected to MongoDB")
}
