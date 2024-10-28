package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
)

func (s *MongoDB) PostFloor(floor models.Floor) error {
	collection := s.Database.Collection("floors")

	_, err := collection.InsertOne(context.TODO(), floor)
	return err
}
