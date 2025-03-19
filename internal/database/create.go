package database

import (
	"fmt"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

// TODO: Add OS to the fields

func Insert(vm VM) (primitive.ObjectID, error) {
	vm.ID = primitive.NewObjectID()

	collection := client.Database("kvmcli").Collection("vms")
	result, err := collection.InsertOne(ctx, vm)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("Insert VM failed:", err)
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not an ObjectID")
	}

	fmt.Println("Inserted VM with _id:", result.InsertedID)
	return insertedID, nil
}
