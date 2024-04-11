package mongo

import (
	"context"
	"time"

	errors "github.com/rotisserie/eris"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitDb() (*mongo.Client, func(), error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, nil, err
	}

	collection := client.Database(viper.GetString("mongo.db")).
		Collection(viper.GetString("mongo.collection"))
	_, err = collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "v_ns", Value: 1}, {Key: "v_name", Value: 1}},
			Options: options.Index().SetName("v_index"),
		},
		{
			Keys:    bson.D{{Key: "u_ns", Value: 1}, {Key: "u_name", Value: 1}},
			Options: options.Index().SetName("u_index"),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to connect to MongoDB")
	}

	var Disconnect = func() {
		client.Disconnect(ctx)
	}
	return client, Disconnect, nil
}
