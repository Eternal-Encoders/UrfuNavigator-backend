package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

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

func (s *MongoDB) GetInstituteIconsByName(names []string) ([]models.InstituteIcon, error) {
	collection := s.Database.Collection("media")

	filter := bson.M{
		"filename": bson.M{
			"$in": names,
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
	log.Println(result[0].Id)

	return result, nil
}

func (s *MongoDB) PostInstituteIcon(icon models.InstituteIconGet) error {
	collection := s.Database.Collection("media")

	_, err := collection.InsertOne(context.TODO(), icon)
	return err
}

func (s *MongoDB) DeleteInstituteIcon(id string) (string, error) {
	collection := s.Database.Collection("media")

	objId, err := primitive.ObjectIDFromHex(id)
	log.Println(objId)
	if err != nil {
		return "", err
	}
	filter := bson.M{
		"_id": objId,
	}

	var icon models.InstituteIcon
	err = collection.FindOne(context.TODO(), filter).Decode(&icon)
	if err != nil {
		return "", err
	}

	_, err = collection.DeleteOne(context.TODO(), filter)
	return icon.Alt, err
}
