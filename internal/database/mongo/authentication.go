package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MongoDB) Register(body models.UserDB) models.ResponseType {
	collection := s.Database.Collection("users")

	filter := bson.M{
		"url": body.Email,
	}

	err := collection.FindOne(context.TODO(), filter).Err()
	if err == nil {
		if err == mongo.ErrNoDocuments {
			return models.ResponseType{Type: 406, Error: errors.New("user with specified email already exists")}
		}
	} else {
		if err != mongo.ErrNoDocuments {
			log.Println(err)
			return models.ResponseType{Type: 500, Error: err}
		}
	}

	// experiment
	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: body.Email}},
			Options: options.Index().SetUnique(true),
		},
	)

	_, err = collection.InsertOne(context.TODO(), body)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

// func (s *MongoDB) Logout() models.ResponseType {
// 	collection := s.Database.Collection("users")

//  return models.ResponseType{Type: 200, Error: nil}
// }
