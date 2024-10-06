package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *MongoDB) GetInstituteIcons(ids []string) ([]models.InstituteIcon, error) {
	collection := s.Database.Collection("media")

	objIds := []primitive.ObjectID{}
	for _, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
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
		return nil, err
	}

	defer curs.Close(context.TODO())

	var result []models.InstituteIcon
	decodeErr := curs.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) GetAllInstituteIcons() ([]models.InstituteIcon, error) {
	collection := s.Database.Collection("media")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var result []models.InstituteIcon
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) PostInstituteIcon(icon models.InstituteIconRequest) error {
	collection := s.Database.Collection("media")

	_, err := collection.InsertOne(context.TODO(), icon)
	return err
}
