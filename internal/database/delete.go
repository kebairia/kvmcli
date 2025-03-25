package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func Delete(name string) error {
	// Create a filter matching the record with the specified name
	filter := bson.M{"name": name}
	collection := client.Database("kvmcli").Collection("vms")
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no record found with name: %s", name)
	}
	return nil
}
