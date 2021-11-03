package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/logger"
)

var MG *mongo.Database

func init() {
	uri := fmt.Sprintf(
		"mongodb://%s:%s",
		config.Env.MongoHost,
		config.Env.MongoPort)
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}
	logger.Sugar.Debugw("Connected to MongoDB!",
		"mongodb config", uri)
	MG = client.Database(config.Env.MongoDBName)
}
