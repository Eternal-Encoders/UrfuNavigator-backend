package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *MongoDB) GetInstitute(url string) (models.Institute, error) {
	collection := s.Database.Collection("institutes")
	filter := bson.M{
		"url": url,
	}

	var result models.Institute
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (s *MongoDB) GetAllInstitutes() ([]models.Institute, error) {
	collection := s.Database.Collection("institutes")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var result []models.Institute
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) PostInstitute(institute models.InstituteRequest) error {
	collection := s.Database.Collection("institutes")

	iconCol := s.Database.Collection("media")
	filter := bson.D{{"alt", institute.Icon}}
	log.Println(filter)
	// var result models.InstituteIcon
	_, err := iconCol.Find(context.TODO(), filter)
	if err != nil {
		log.Println("1")
		return err
	}

	_, err = collection.InsertOne(context.TODO(), institute)
	log.Println("2")
	return err
}
