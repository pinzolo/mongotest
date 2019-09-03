package mongotest

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect(ctx context.Context) (context.Context, *mongo.Client, context.CancelFunc, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(conf.url))
	if err != nil {
		return ctx, nil, nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(conf.timeout)*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		cancel()
		return ctx, nil, nil, err
	}
	return ctx, client, func() {
		client.Disconnect(ctx)
		cancel()
	}, nil
}

func connectCollection(ctx context.Context, collName string) (context.Context, *mongo.Collection, context.CancelFunc, error) {
	ctx, client, cancel, err := connect(ctx)
	if err != nil {
		return ctx, nil, cancel, err
	}
	collection := client.Database(conf.database).Collection(collName)
	return ctx, collection, cancel, err
}
