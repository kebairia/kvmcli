package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// GetObjectsByNamespace retrieves all documents of type T from the specified collection
// that match the given namespace.
func GetObjectsByNamespace[T VMRecord | NetRecord | StoreRecord](
	namespace, collectionName string,
) ([]T, error) {
	collection := client.Database(Database).Collection(collectionName)
	filter := bson.M{"namespace": namespace}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf(
			"error finding documents for namespace %s in collection %s: %w",
			namespace,
			collectionName,
			err,
		)
	}
	defer cursor.Close(ctx)

	var objects []T
	for cursor.Next(ctx) {
		var obj T
		if err := cursor.Decode(&obj); err != nil {
			return nil, fmt.Errorf("error decoding record: %w", err)
		}
		objects = append(objects, obj)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return objects, nil
}

// Get retrieves a single VMRecord by its name.
func GetRecord[T VMRecord | NetRecord | StoreRecord](name, collectionName string) (T, error) {
	collection := client.Database("kvmcli").Collection(collectionName)
	// filter := bson.M{"name": name}
	var filter bson.M
	var dummy T
	switch any(dummy).(type) {
	case StoreRecord:
		filter = bson.M{"metadata.name": name}
	default:
		filter = bson.M{"name": name}
	}

	var record T
	err := collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error finding document for name %s: %w", name, err)
	}
	return record, nil
}
