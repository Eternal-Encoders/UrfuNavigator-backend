package mongo

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoDB) GetFloor(context context.Context, filter []models.Query) (models.Floor, models.ResponseType) {
	collection := s.Database.Collection("floors")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.Floor{}, models.ResponseType{Type: 500, Error: err}
	}

	var result models.Floor
	err = collection.FindOne(context, bsonFilter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified filter")}
		} else {
			return result, models.ResponseType{Type: 500, Error: err}
		}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetManyFloors(context context.Context, filter []models.Query) ([]models.Floor, models.ResponseType) {
	collection := s.Database.Collection("floors")
	var result []models.Floor

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return result, models.ResponseType{Type: 500, Error: err}
	}

	cursor, err := collection.Find(context, bsonFilter)
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context)

	decodeErr := cursor.All(context, &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetAllFloors(context context.Context) ([]models.Floor, models.ResponseType) {
	collection := s.Database.Collection("floors")

	cursor, err := collection.Find(context, bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context)

	var result []models.Floor
	decodeErr := cursor.All(context, &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) InsertFloor(context context.Context, floor models.FloorRequest) models.ResponseType {
	collection := s.Database.Collection("floors")

	_, err := collection.InsertOne(context, floor)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) DeleteFloor(context context.Context, filter []models.Query) models.ResponseType {
	collection := s.Database.Collection("floors")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, res := s.GetFloor(context, filter)
	if res.Error != nil {
		return res
	}

	_, err = collection.DeleteOne(context, bsonFilter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) DeleteManyFloors(context context.Context, filter []models.Query) models.ResponseType {
	collection := s.Database.Collection("floors")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, err = collection.DeleteMany(context, bsonFilter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) UpdateFloor(context context.Context, filter []models.Query, body models.FloorPut) models.ResponseType {
	collection := s.Database.Collection("floors")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, res := s.GetFloor(context, filter)
	if res.Error != nil {
		return res
	}

	_, err = collection.UpdateOne(context, bsonFilter, bson.M{"$set": body})
	fmt.Println(err != nil)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) UpdateManyFloors(context context.Context, filter []models.Query, body interface{}) models.ResponseType {
	collection := s.Database.Collection("floors")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	updateBody := utils.ToBson(body)

	_, err = collection.UpdateMany(context, bsonFilter, bson.M{"$set": updateBody})
	fmt.Println(err != nil)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}
