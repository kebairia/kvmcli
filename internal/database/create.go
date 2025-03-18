package database

import (
	"fmt"
	"log"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

// TODO: Add OS to the fields

func CreateVMEntry(
	name string,
	namespace string,
	ram int,
	cpu int,
	macaddress string,
	network string,
) error {
	vm := VM{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Namespace:  namespace,
		CPU:        cpu,
		RAM:        ram,
		MacAddress: macaddress,
		NetworkID:  network,
	}
	collection := client.Database("kvmcli").Collection("vms")
	result, err := collection.InsertOne(ctx, vm)
	if err != nil {
		log.Fatal("Insert VM failed:", err)
	}

	fmt.Println("Inserted VM with _id:", result.InsertedID)
	return nil
}
