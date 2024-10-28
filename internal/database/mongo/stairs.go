package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *MongoDB) PostStairs(graphs []*models.GraphPoint) error {
	collection := s.Database.Collection("stairs")

	if len(graphs) == 0 {
		return nil
	}

	for i := range graphs {
		if graphs[i].StairId != "" {
			filter := bson.D{{"_id", graphs[i].StairId}}
			var stairGraph models.Stair
			err := collection.FindOne(context.TODO(), filter).Decode(&stairGraph)
			if err != nil {
				_, err = collection.InsertOne(context.TODO(), models.Stair{
					Id:         graphs[i].StairId,
					StairPoint: graphs[i].StairId,
					Institute:  graphs[i].Institute,
					Links:      []string{graphs[i].Id},
				})

				if err != nil {
					return err
				}
			} else {
				l := append(stairGraph.Links, graphs[i].Id)
				_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"links": l}})
			}
		}
	}

	return nil
}
