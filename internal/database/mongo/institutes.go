package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *MongoDB) GetInstitute(url string) models.InstituteReadDBResponse {
	collection := s.Database.Collection("institutes")
	filter := bson.M{
		"url": url,
	}

	var result models.Institute
	var errType int
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		errType = fiber.StatusNotFound
	} else {
		errType = fiber.StatusInternalServerError
	}

	return models.InstituteReadDBResponse{
		Response:  result,
		Error:     err,
		ErrorType: errType,
	}
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
	// log.Println(result[0].Id)

	return result, nil
}

func (s *MongoDB) PostInstitute(institute models.InstituteRequest) error {
	collection := s.Database.Collection("institutes")
	iconCol := s.Database.Collection("media")

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		filter := bson.M{"alt": institute.Icon}

		err := iconCol.FindOne(ctx, filter).Err()
		if err != nil {
			log.Println("Icon error:", err)
			return nil, err
		}

		result, err := collection.InsertOne(ctx, institute)
		return result, err
	}, txnOptions)

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
