package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	log.Println(result[0].Id)

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

func (s *MongoDB) DeleteInstitute(id string) error {
	collection := s.Database.Collection("institutes")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": objId,
	}

	_, err = collection.DeleteOne(context.TODO(), filter)
	return err
}

func (s *MongoDB) UpdateInstitute(body models.InstituteRequest, id string) error {
	collection := s.Database.Collection("institutes")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err != nil)
		return err
	}
	filter := bson.M{
		"_id": objId,
	}

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": body})
	fmt.Println(err != nil)
	return err
}
