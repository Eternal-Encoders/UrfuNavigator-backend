package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
)

func (s *MongoDB) PostGraphs(graphs []*models.GraphPoint) error {
	collection := s.Database.Collection("graph_points")

	if len(graphs) == 0 {
		return nil
	}

	newValue := make([]interface{}, len(graphs))

	for i := range graphs {
		newValue[i] = graphs[i]
	}

	_, err := collection.InsertMany(context.TODO(), newValue)
	return err
}
