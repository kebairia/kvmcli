package database

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertVMold(record *VMRecord) (primitive.ObjectID, error) {
	record.ID = primitive.NewObjectID()
	collection := client.Database("kvmcli").Collection(VMsCollection)
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

func InsertNet(record *NetRecord) (primitive.ObjectID, error) {
	record.ID = primitive.NewObjectID()
	collection := client.Database("kvmcli").Collection(NetworksCollection)
	result, err := collection.InsertOne(ctx, record)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("Insert record failed: %w", err)
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not an ObjectID")
	}

	logger.Log.Debugf("Inserted Network record with _id: %v", result.InsertedID)
	return insertedID, nil
}

func InsertStore(record *StoreRecord) (primitive.ObjectID, error) {
	record.ID = primitive.NewObjectID()
	collection := client.Database("kvmcli").Collection(StoreCollection)
	result, err := collection.InsertOne(ctx, record)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("Insert record failed: %w", err)
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("InsertedID is not an ObjectID")
	}

	logger.Log.Debugf("Inserted store record with _id: %v", result.InsertedID)
	return insertedID, nil
}
