package database

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertVM(record *VMRecord) (primitive.ObjectID, error) {
	record.ID = primitive.NewObjectID()
	collection := client.Database("kvmcli").Collection("vms")
	result, err := collection.InsertOne(ctx, record)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("Insert record failed: %w", err)
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not an ObjectID")
	}

	logger.Log.Debugf("Inserted VM record with _id: %v", result.InsertedID)
	return insertedID, nil
}
