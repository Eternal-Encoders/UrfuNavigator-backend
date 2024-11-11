package database

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *MongoDB) PostFloor(floor models.FloorFromFile) error {
	floorCol := s.Database.Collection("floors")
	graphCol := s.Database.Collection("graph_points")
	stairsCol := s.Database.Collection("stairs")
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return err
	}

	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		audArr := []*models.Auditorium{}
		for _, v := range floor.Audiences {
			audArr = append(audArr, v)
		}

		graphArr := []*models.GraphPoint{}
		graphKeysArr := []string{}
		for k, v := range floor.Graph {
			graphArr = append(graphArr, v)
			graphKeysArr = append(graphKeysArr, k)
		}

		floorReq := models.FloorRequest{
			Institute: floor.Institute,
			Floor:     floor.Floor,
			Width:     floor.Width,
			Height:    floor.Height,
			Service:   floor.Service,
			Audiences: audArr,
			Graph:     graphKeysArr,
		}

		_, err = floorCol.InsertOne(ctx, floorReq)
		if err != nil {
			log.Println("Something went wrong with inserting floor obj in DB")
			return nil, err
		}

		if len(graphArr) == 0 {
			log.Println("Empty graphs list")
			return nil, errors.New("empty graphs list")
		}

		newValue := make([]interface{}, len(graphArr))

		for i := range graphArr {
			newValue[i] = graphArr[i]
		}

		_, err := graphCol.InsertMany(ctx, newValue)
		if err != nil {
			log.Println("Something went wrong with inserting graph_point objects in DB")
			return nil, err
		} else {
			for i := range graphArr {
				if graphArr[i].StairId != "" {
					filter := bson.M{"_id": graphArr[i].StairId}
					var stairGraph models.Stair

					err := stairsCol.FindOne(ctx, filter).Decode(&stairGraph)
					if err != nil {
						_, err = stairsCol.InsertOne(ctx, models.Stair{
							Id:         graphArr[i].StairId,
							StairPoint: graphArr[i].StairId,
							Institute:  graphArr[i].Institute,
							Links:      []string{graphArr[i].Id},
						})

						if err != nil {
							log.Println("Something went wrong with inserting stair object in DB")
							return nil, err
						}
					} else {
						l := append(stairGraph.Links, graphArr[i].Id)
						_, err = stairsCol.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"links": l}})
						if err != nil {
							log.Println("Something went wrong in stairsCol.UpdateOne")
							return nil, err
						}
					}
				}
			}
		}

		return nil, nil
	}, txnOptions)

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
