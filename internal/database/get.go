package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// ISSUE: This retrieve only VMRecord, therefore I need another function to retreive
// snapshots, network information.
// FIX: I need to create an interface for database for all types of records

// GetVMsByNamespace retrieves all VMRecord documents that match the given namespace.
func GetVMsByNamespace(namespace string) ([]VMRecord, error) {
	collection := client.Database("kvmcli").Collection("vms")
	filter := bson.M{"namespace": namespace}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding documents for namespace %s: %w", namespace, err)
	}
	defer cursor.Close(ctx)

	var vms []VMRecord
	for cursor.Next(ctx) {
		var vm VMRecord
		if err := cursor.Decode(&vm); err != nil {
			return nil, fmt.Errorf("error decoding VM record: %w", err)
		}
		vms = append(vms, vm)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return vms, nil
}

// Get retrieves a single VMRecord by its name.
func GetVM(name string) (VMRecord, error) {
	collection := client.Database("kvmcli").Collection("vms")
	filter := bson.M{"name": name}

	var vm VMRecord
	err := collection.FindOne(ctx, filter).Decode(&vm)
	if err != nil {
		return VMRecord{}, fmt.Errorf("error finding document for name %s: %w", name, err)
	}
	return vm, nil
}

// Get retrieves a single VMRecord by its name.
func GetNetwork(name string) (NetRecord, error) {
	collection := client.Database("kvmcli").Collection("networks")
	filter := bson.M{"name": name}

	var network NetRecord
	err := collection.FindOne(ctx, filter).Decode(&network)
	if err != nil {
		return NetRecord{}, fmt.Errorf("error finding document for name %s: %w", name, err)
	}
	return network, nil
}
