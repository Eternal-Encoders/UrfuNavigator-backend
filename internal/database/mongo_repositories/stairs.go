package mongo

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoDB) GetStair(context context.Context, filter []models.Query) (models.Stair, models.ResponseType) {
	collection := s.Database.Collection("stairs")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.Stair{}, models.ResponseType{Type: 500, Error: err}
	}

	var result models.Stair
	err = collection.FindOne(context, bsonFilter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified filter")}
		} else {
			return result, models.ResponseType{Type: 500, Error: err}
		}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetManyStairs(context context.Context, filter []models.Query) ([]models.Stair, models.ResponseType) {
	collection := s.Database.Collection("stairs")
	var result []models.Stair

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

func (s *MongoDB) GetAllStairs(context context.Context) ([]models.Stair, models.ResponseType) {
	collection := s.Database.Collection("stairs")

	cursor, err := collection.Find(context, bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context)

	var result []models.Stair
	decodeErr := cursor.All(context, &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) InsertManyStairs(context context.Context, graphs []*models.Stair) models.ResponseType {
	collection := s.Database.Collection("stairs")

	newValue := make([]interface{}, len(graphs))

	for i := range graphs {
		newValue[i] = graphs[i]
	}

	if len(newValue) > 0 {
		_, err := collection.InsertMany(context, newValue)
		if err != nil {
			log.Println("graphs insertion error")
			return models.ResponseType{Type: 500, Error: err}
		}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) DeleteStair(context context.Context, filter []models.Query) models.ResponseType {
	collection := s.Database.Collection("stairs")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, res := s.GetStair(context, filter)
	if res.Error != nil {
		return res
	}

	_, err = collection.DeleteOne(context, bsonFilter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) DeleteManyStairs(context context.Context, filter []models.Query) models.ResponseType {
	collection := s.Database.Collection("stairs")

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

func (s *MongoDB) UpdateStair(context context.Context, filter []models.Query, body models.Stair) models.ResponseType {
	collection := s.Database.Collection("stairs")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, res := s.GetStair(context, filter)
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

func (s *MongoDB) UpdateManyStairs(context context.Context, filter []models.Query, body interface{}) models.ResponseType {
	collection := s.Database.Collection("stairs")

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
