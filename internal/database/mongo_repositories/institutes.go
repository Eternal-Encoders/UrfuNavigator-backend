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

func (s *MongoDB) GetInstitute(context context.Context, filter []models.Query) (models.Institute, models.ResponseType) {
	collection := s.Database.Collection("insitutes")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.Institute{}, models.ResponseType{Type: 500, Error: err}
	}

	var result models.Institute

	err = collection.FindOne(context, bsonFilter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified filter")}
		} else {
			return result, models.ResponseType{Type: 500, Error: err}
		}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetManyInstitutes(context context.Context, filter []models.Query) ([]models.Institute, models.ResponseType) {
	collection := s.Database.Collection("insitutes")
	var result []models.Institute

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
	// if len(result) == 0 {
	// 	return nil, models.ResponseType{Type: 404, Error: errors.New("there are no institutes")}
	// }

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetAllInstitutes(context context.Context) ([]models.Institute, models.ResponseType) {
	collection := s.Database.Collection("insitutes")

	cursor, err := collection.Find(context, bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context)

	var result []models.Institute
	decodeErr := cursor.All(context, &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}
	// if len(result) == 0 {
	// 	return nil, models.ResponseType{Type: 404, Error: errors.New("there are no institutes")}
	// }

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) InsertInstitute(context context.Context, institute models.InstitutePost) models.ResponseType {
	collection := s.Database.Collection("insitutes")

	filter := bson.M{"name": institute.Name}

	err := collection.FindOne(context, filter).Err()
	if err == nil {
		return models.ResponseType{Type: 406, Error: errors.New("institute already exists")}
	}

	_, err = collection.InsertOne(context, institute)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

// func (s *MongoDB) InsertManyInstitutes(context context.Context, institutes []models.InstitutePost) models.ResponseType {

// }

func (s *MongoDB) DeleteInstitute(context context.Context, filter []models.Query) models.ResponseType {
	collection := s.Database.Collection("insitutes")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, res := s.GetInstitute(context, filter)
	if res.Error != nil {
		return res
	}

	_, err = collection.DeleteOne(context, bsonFilter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) DeleteManyInstitutes(context context.Context, filter []models.Query) models.ResponseType {
	collection := s.Database.Collection("insitutes")

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

func (s *MongoDB) UpdateInstitute(context context.Context, filter []models.Query, body models.InstitutePost) models.ResponseType {
	collection := s.Database.Collection("insitutes")

	bsonFilter, err := utils.CreateBSONFilter(filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	_, res := s.GetInstitute(context, filter)
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

func (s *MongoDB) UpdateManyInstitutes(context context.Context, filter []models.Query, body interface{}) models.ResponseType {
	collection := s.Database.Collection("insitutes")

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
