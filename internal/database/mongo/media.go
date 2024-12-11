package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *MongoDB) GetInstituteIcons(ids []string) ([]models.InstituteIcon, models.ResponseType) {
	collection := s.Database.Collection("media")

	objIds := []primitive.ObjectID{}
	for _, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, models.ResponseType{Type: 500, Error: err}
		}
		objIds = append(objIds, objId)
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	}
	curs, err := collection.Find(context.TODO(), filter)
	if err != nil {
		if err == models.ErrNoDocuments {

		} else {
			return nil, models.ResponseType{Type: 500, Error: err}
		}
	}

	defer curs.Close(context.TODO())

	var result []models.InstituteIcon
	decodeErr := curs.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: decodeErr}
	}

	// if len(result) == 0 {
	// 	return nil, models.ResponseType{Type: 404, Error: errors.New("there are no icons with specified ids")}
	// }

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetInstituteIconsByName(names []string) ([]models.InstituteIcon, models.ResponseType) {
	collection := s.Database.Collection("media")

	filter := bson.M{
		"filename": bson.M{
			"$in": names,
		},
	}
	curs, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer curs.Close(context.TODO())

	var result []models.InstituteIcon
	decodeErr := curs.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: decodeErr}
	}

	// if len(result) == 0 {
	// 	return nil, models.ResponseType{Type: 404, Error: errors.New("there are no icons with specified names")}
	// }

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetAllInstituteIcons() ([]models.InstituteIcon, models.ResponseType) {
	collection := s.Database.Collection("media")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context.TODO())

	var result []models.InstituteIcon
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: decodeErr}
	}
	log.Println(result[0].Id)

	// if len(result) == 0 {
	// 	return nil, models.ResponseType{Type: 404, Error: errors.New("there are no icons")}
	// }

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) PostInstituteIcon(icon models.InstituteIconGet) models.ResponseType {
	collection := s.Database.Collection("media")

	_, err := collection.InsertOne(context.TODO(), icon)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) DeleteInstituteIcon(id string) (string, models.ResponseType) {
	collection := s.Database.Collection("media")

	objId, err := primitive.ObjectIDFromHex(id)
	log.Println(objId)
	if err != nil {
		return "", models.ResponseType{Type: 500, Error: err}
	}
	filter := bson.M{
		"_id": objId,
	}

	var icon models.InstituteIcon
	err = collection.FindOne(context.TODO(), filter).Decode(&icon)
	if err != nil {
		if err == models.ErrNoDocuments {
			return "", models.ResponseType{Type: 404, Error: errors.New("there is no icon with specified id: " + id)}
		} else {
			return "", models.ResponseType{Type: 500, Error: err}
		}
	}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return "", models.ResponseType{Type: 500, Error: err}
	}

	return icon.Alt, models.ResponseType{Type: 200, Error: nil}
}
