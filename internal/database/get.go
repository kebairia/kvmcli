package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	database   = "kvmcli"
	collection = "vms"
)

func GetALLVMs() ([]VMRecord, error) {
	collection := client.Database("kvmcli").Collection("vms")
	filter := bson.M{}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
		defer cursor.Close(ctx)
	}
	var results []VMRecord
	for cursor.Next(ctx) {
		var entry VMRecord
		if err := cursor.Decode(&entry); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		results = append(results, entry)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}
	return results, nil
}
