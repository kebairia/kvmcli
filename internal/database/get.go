package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// ISSUE: This retreive only VMRecord, therefore I need another function to retreive
// snapshots, network information.
// FIX: I need to create an interface for database for all types of records

func Get(name string) (infos VMRecord, err error) {
	collection := client.Database("kvmcli").Collection("vms")
	filter := bson.M{
		"name": name,
	}
	var vm VMRecord
	err = collection.FindOne(ctx, filter).Decode(&vm)
	if err != nil {
		return VMRecord{}, fmt.Errorf("Error finding document: %w", err)
	}
	return vm, nil
}
