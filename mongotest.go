package mongotest

import (
	"context"

	"github.com/tkuchiki/parsetime"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Try connecting to MongoDB server.
func Try() error {
	if err := validateConfig(); err != nil {
		return err
	}
	ctx, client, cancel, err := connect(context.Background())
	if err != nil {
		return err
	}
	defer cancel()
	return client.Ping(ctx, nil)
}

// CountWithContext returns document count in collection that has given name.
// This function uses given context.
func CountWithContext(ctx context.Context, collectionName string) (int64, error) {
	ctx, coll, cancel, err := connectCollection(ctx, collectionName)
	if err != nil {
		return 0, err
	}
	defer cancel()
	n, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return n, nil
}

// CountIntWithContext returns document count as int in collection that has given name.
// This function uses given context.
func CountIntWithContext(ctx context.Context, collectionName string) (int, error) {
	n, err := CountWithContext(ctx, collectionName)
	return int(n), err
}

// Count returns document count in collection that has given name.
func Count(collectionName string) (int64, error) {
	return CountWithContext(context.Background(), collectionName)
}

// CountInt returns document count as int in collection that has given name.
func CountInt(collectionName string) (int, error) {
	return CountIntWithContext(context.Background(), collectionName)
}

// FindWithContext document that has given id in given named collection.
// This function uses given context.
func FindWithContext(ctx context.Context, collectionName string, id interface{}) (map[string]interface{}, error) {
	ctx, coll, cancel, err := connectCollection(ctx, collectionName)
	if err != nil {
		return nil, err
	}
	defer cancel()
	var doc bson.M
	err = coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	return doc, err
}

// Find document that has given id in given named collection.
func Find(collectionName string, id interface{}) (map[string]interface{}, error) {
	return FindWithContext(context.Background(), collectionName, id)
}

// SimpleConvertTime provides simple PreInsertFunc for converting string time to time.Time.
func SimpleConvertTime(collectionName, fieldName string) PreInsertFunc {
	return func(collName string, value DocData) (DocData, error) {
		if collName == collectionName {
			sv, ok := value.StringValue(fieldName)
			if !ok {
				return value, nil
			}
			p, err := parsetime.NewParseTime()
			if err != nil {
				return nil, err
			}
			t, err := p.Parse(sv)
			if err != nil {
				return nil, err
			}
			value[fieldName] = t
		}
		return value, nil
	}
}

// SimpleConvertBytes provides simple PreInsertFunc for converting string to BinData
func SimpleConvertBytes(collectionName, fieldName string) PreInsertFunc {
	return func(collName string, value DocData) (DocData, error) {
		if collName == collectionName {
			sv, ok := value.StringValue(fieldName)
			if !ok {
				return value, nil
			}
			value[fieldName] = []byte(sv)
		}
		return value, nil
	}
}

// SimpleConvertObjID provides simple PreInsertFunc for converting string to ObjectID
func SimpleConvertObjID(collectionName, fieldName string) PreInsertFunc {
	return func(collName string, value DocData) (DocData, error) {
		if collName == collectionName {
			sv, ok := value.StringValue(fieldName)
			if !ok {
				return value, nil
			}
			objID, err := primitive.ObjectIDFromHex(sv)
			if err != nil {
				return nil, err
			}
			value[fieldName] = objID
		}
		return value, nil
	}
}

// DropCollections drops multiple collections with the given names.
func DropCollections(collections ...string) error {
	for _, collection := range collections {
		ctx := context.Background()
		ctx, coll, cancel, err := connectCollection(ctx, collection)
		if err != nil {
			return err
		}
		defer cancel()
		if err := coll.Drop(ctx); err != nil {
			return err
		}
	}
	return nil
}

// DropDatabase drops the entire database.
func DropDatabase() error {
	ctx := context.Background()
	ctx, client, cancel, err := connect(ctx)
	if err != nil {
		return err
	}
	defer cancel()
	return client.Database(conf.Database).Drop(ctx)
}
