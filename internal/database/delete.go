package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func Remove(name string) error {
	collection := client.Database("kvmcli").Collection("vms")
	filter := bson.M{"name": name}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to remove VM entry from database: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no database entry found for VM: %s", name)
	}
	return nil
}
