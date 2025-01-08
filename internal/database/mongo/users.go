package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoDB) GetUser(email string) (models.UserDB, models.ResponseType) {
	collection := s.Database.Collection("users")
	filter := bson.M{
		"email": email,
	}

	var result models.UserDB
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, models.ResponseType{Type: 404, Error: errors.New("there is no user with specified email: " + email)}
		} else {
			return result, models.ResponseType{Type: 500, Error: err}
		}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetAllUsers() ([]models.UserDB, models.ResponseType) {
	collection := s.Database.Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context.TODO())

	var result []models.UserDB
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}
