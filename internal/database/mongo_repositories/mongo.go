package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(uri string, collection string) *MongoDB {
	dbOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), dbOptions)

	if err != nil {
		log.Fatal(err)
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(collection),
	}
}

func (s *MongoDB) Disconnect() error {
	return s.Client.Disconnect(context.TODO())
}
