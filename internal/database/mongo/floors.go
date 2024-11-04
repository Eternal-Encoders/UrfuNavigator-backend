package database

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *MongoDB) PostFloor(floor models.FloorRequest) error {
	collection := s.Database.Collection("floors")

	_, err := collection.InsertOne(context.TODO(), floor)
	return err
}

func (s *MongoDB) GetFloor(id string) (models.Floor, error) {
	collection := s.Database.Collection("floors")
	log.Println(id)
	objId, err := primitive.ObjectIDFromHex(id)
	log.Println(objId)
	if err != nil {
		return models.Floor{}, err
	}

	filter := bson.M{
		"_id": objId,
	}

	var result models.Floor
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func (s *MongoDB) GetAllFloors() ([]models.Floor, error) {
	collection := s.Database.Collection("floors")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var result []models.Floor
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) DeleteFloor(id string) error {
	stairsCol := s.Database.Collection("stairs")
	graphsCol := s.Database.Collection("graph_points")
	floorsCol := s.Database.Collection("floors")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	floorFilter := bson.M{
		"_id": objId,
	}
	var floor models.Floor
	err = floorsCol.FindOne(context.TODO(), floorFilter).Decode(&floor)
	if err != nil {
		return err
	}

	for _, v := range floor.Graph {
		graphFilter := bson.M{
			"_id": v,
		}
		var graph models.GraphPoint
		err := graphsCol.FindOne(context.TODO(), graphFilter).Decode(&graph)
		if err != nil {
			return err
		}

		if graph.StairId != "" {
			stairFilter := bson.M{
				"stairPoint": graph.StairId,
			}
			var stair models.Stair
			err = stairsCol.FindOne(context.TODO(), stairFilter).Decode(&stair)
			if err != nil {
				return err
			}
			log.Println(stair.Id)

			linkIndex := utils.GetIndex(stair.Links, graph.Id)
			newLinks := append(stair.Links[:linkIndex], stair.Links[linkIndex+1:]...)
			_, err = stairsCol.UpdateOne(context.TODO(), stairFilter, bson.M{"$set": bson.M{"links": newLinks}})
			if err != nil {
				return err
			}
		}

		_, err = graphsCol.DeleteOne(context.TODO(), graphFilter)
		if err != nil {
			return err
		}
	}

	_, err = floorsCol.DeleteOne(context.TODO(), floorFilter)
	if err != nil {
		return err
	}

	return err
}
