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

func (s *MongoDB) GetGraph(context context.Context, filter []models.Query) (models.GraphPoint, models.ResponseType) {
	collection := s.Database.Collection("graph_points")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.GraphPoint{}, models.ResponseType{Type: 500, Error: err}
	}

	var result models.GraphPoint
	err = collection.FindOne(context, bsonFilter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified filter")}
		} else {
			return result, models.ResponseType{Type: 500, Error: err}
		}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetManyGraphs(context context.Context, filter []models.Query) ([]models.GraphPoint, models.ResponseType) {
	collection := s.Database.Collection("graph_points")
	var result []models.GraphPoint

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

func (s *MongoDB) GetManyGraphsByIds(context context.Context, ids []string) ([]models.GraphPoint, models.ResponseType) {
	collection := s.Database.Collection("graph_points")
	var result []models.GraphPoint

	// bsonFilter, err := utils.CreateBSONFilter(filter)
	// if err != nil {
	// 	return result, models.ResponseType{Type: 500, Error: err}
	// }

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	cursor, err := collection.Find(context, filter)
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

func (s *MongoDB) GetAllGraphs(context context.Context) ([]models.GraphPoint, models.ResponseType) {
	collection := s.Database.Collection("graph_points")

	cursor, err := collection.Find(context, bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context)

	var result []models.GraphPoint
	decodeErr := cursor.All(context, &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) InsertManyGraphs(context context.Context, graphs []*models.GraphPoint) models.ResponseType {
	collection := s.Database.Collection("graph_points")

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

func (s *MongoDB) DeleteGraph(context context.Context, filter []models.Query) models.ResponseType {
	collection := s.Database.Collection("graph_points")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, res := s.GetGraph(context, filter)
	if res.Error != nil {
		return res
	}

	_, err = collection.DeleteOne(context, bsonFilter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) DeleteManyGraphs(context context.Context, filter []models.Query) models.ResponseType {
	collection := s.Database.Collection("graph_points")

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

func (s *MongoDB) DeleteManyGraphsByIds(context context.Context, ids []string) models.ResponseType {
	collection := s.Database.Collection("graph_points")

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	_, err := collection.DeleteMany(context, filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) UpdateGraph(context context.Context, filter []models.Query, body models.GraphPointPut) models.ResponseType {
	collection := s.Database.Collection("graph_points")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, res := s.GetGraph(context, filter)
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

func (s *MongoDB) UpdateManyGraphs(context context.Context, filter []models.Query, body interface{}) models.ResponseType {
	collection := s.Database.Collection("graph_points")

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
