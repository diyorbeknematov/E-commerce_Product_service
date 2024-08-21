package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() (*mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().
		ApplyURI("mongodb://mongo:27017"))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database("e-commerce"), nil
}
