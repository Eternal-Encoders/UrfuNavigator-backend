package database

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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

func (s *MongoDB) GetGraph(id string) (models.GraphPoint, error) {
	collection := s.Database.Collection("graph_points")

	filter := bson.M{
		"_id": id,
	}

	var result models.GraphPoint
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func (s *MongoDB) GetAllGraphs() ([]models.GraphPoint, error) {
	collection := s.Database.Collection("graph_points")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var result []models.GraphPoint
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) DeleteGraph(id string) error {
	stairsCol := s.Database.Collection("stairs")
	graphsCol := s.Database.Collection("graph_points")
	floorsCol := s.Database.Collection("floors")
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return err
	}

	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {

		graphFilter := bson.M{
			"_id": id,
		}
		var graph models.GraphPoint
		err := graphsCol.FindOne(context.TODO(), graphFilter).Decode(&graph)
		if err != nil {
			return nil, err
		}

		floorFilter := bson.M{
			"floor":     graph.Floor,
			"institute": graph.Institute,
		}
		var floor models.Floor
		err = floorsCol.FindOne(context.TODO(), floorFilter).Decode(&floor)
		if err != nil {
			return nil, err
		}
		log.Println(floor.Institute, floor.Floor)

		if graph.StairId != "" {
			stairFilter := bson.M{
				"stairPoint": graph.StairId,
			}
			var stair models.Stair
			err = stairsCol.FindOne(context.TODO(), stairFilter).Decode(&stair)
			if err != nil {
				return nil, err
			}
			log.Println(stair.Id)

			linkIndex := utils.GetIndex(stair.Links, graph.Id)
			newLinks := append(stair.Links[:linkIndex], stair.Links[linkIndex+1:]...)
			_, err = stairsCol.UpdateOne(context.TODO(), stairFilter, bson.M{"$set": bson.M{"links": newLinks}})
			if err != nil {
				return nil, err
			}
		}

		graphIndex := utils.GetIndex(floor.Graph, graph.Id)
		newGraphs := append(floor.Graph[:graphIndex], floor.Graph[graphIndex+1:]...)
		_, err = floorsCol.UpdateOne(context.TODO(), bson.M{"_id": floor.Id}, bson.M{"$set": bson.M{"graph": newGraphs}})
		if err != nil {
			return nil, err
		}

		_, err = graphsCol.DeleteOne(context.TODO(), graphFilter)
		if err != nil {
			return nil, err
		}

		return nil, err

	}, txnOptions)

	return err
}
