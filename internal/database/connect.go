package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo() {
	uri := "mongodb://root:example@localhost:27017/?authSource=admin"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}
	// Define a context with a timeout for the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Connect to MongoDB.
	if err := client.Connect(ctx); err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	// Ensure the client disconnects when the function ends.
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
	}()
	db := client.Database("kvmcli")
	collection := db.Collection("vms")
	// Create a new VM document.
	vm := bson.D{
		{"name", "vm1"},
		{"status", "running"},
		{"cpu", 2},
		{"memory", 4096},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
		// You can add more fields like disks, network details, etc.
	}

	// Insert the document into the collection.
	insertResult, err := collection.InsertOne(ctx, vm)
	if err != nil {
		log.Fatal("Error inserting document:", err)
	}
	fmt.Println("Inserted VM with ID:", insertResult.InsertedID)

	// Query the inserted VM document by its name.
	var result bson.M
	filter := bson.D{{"name", "vm1"}}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal("Error finding document:", err)
	}
	fmt.Println("Found VM document:", result)
}
