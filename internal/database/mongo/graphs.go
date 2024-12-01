package database

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"errors"
	"log"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *MongoDB) PostGraphs(context context.Context, graphs []*models.GraphPoint) error {
	collection := s.Database.Collection("graph_points")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return err
	}

	defer session.EndSession(context)

	_, err = session.WithTransaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		newValue := make([]interface{}, len(graphs))

		for i := range graphs {
			newValue[i] = graphs[i]
		}

		if len(newValue) > 0 {
			_, err := collection.InsertMany(ctx, newValue)
			if err != nil {
				log.Println("graphs insertion error")
				return nil, err
			}
		}

		err = s.PostStairs(ctx, graphs)
		if err != nil {
			log.Println("stair error")
			return nil, err
		}

		return nil, nil
	}, txnOptions)

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

func (s *MongoDB) UpdateGraph(context context.Context, body models.GraphPointPut, id string) error {
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

	defer session.EndSession(context)

	_, err = session.WithTransaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		graphFilter := bson.M{
			"_id": id,
		}
		var graph models.GraphPoint
		err := graphsCol.FindOne(ctx, graphFilter).Decode(&graph)
		if err != nil {
			return nil, err
		}

		//Floors and institute

		var oldFloor models.Floor
		if graph.Institute != "" {
			oldFloorFilter := bson.M{
				"floor":     graph.Floor,
				"institute": graph.Institute,
			}
			err = floorsCol.FindOne(ctx, oldFloorFilter).Decode(&oldFloor)
			if err != nil {
				return nil, err
			}
		}

		var newFloor models.Floor
		if body.Institute != "" {
			newFloorFilter := bson.M{
				"floor":     body.Floor,
				"institute": body.Institute,
			}
			err = floorsCol.FindOne(ctx, newFloorFilter).Decode(&newFloor)
			if err != nil {
				return nil, err
			}
		}

		if newFloor.Id != oldFloor.Id {
			if body.Institute != "" {
				newFloor.Graph = append(newFloor.Graph, id)

				_, err = floorsCol.UpdateOne(ctx, bson.M{"_id": newFloor.Id}, bson.M{"$set": bson.M{"graph": newFloor.Graph}})
				if err != nil {
					return nil, err
				}
			}

			if graph.Institute != "" {
				i := utils.GetIndex(oldFloor.Graph, id)
				if i != -1 {
					oldFloor.Graph = append(oldFloor.Graph[:i], oldFloor.Graph[i+1:]...)
				}

				for _, aud := range oldFloor.Audiences {
					if aud.PointId == id {
						aud.PointId = ""
						//Если для каждой аудитории уникальная точка
						break
					}
				}

				_, err = floorsCol.UpdateOne(ctx, bson.M{"_id": oldFloor.Id}, bson.M{"$set": oldFloor})
				if err != nil {
					return nil, err
				}
			}
		}

		added, deleted := utils.GetAddedDeleted(graph.Links, body.Links)

		if len(deleted) > 0 {
			var deletedGraphs []models.GraphPoint
			delFilter := bson.M{
				"_id": bson.M{
					"$in": deleted,
				},
			}

			cur, err := graphsCol.Find(ctx, delFilter)
			if err != nil {
				return nil, err
			}

			err = cur.All(ctx, &deletedGraphs)
			if err != nil {
				return nil, err
			}

			cur.Close(ctx)

			for _, graphLink := range deletedGraphs {
				links := graphLink.Links
				i := utils.GetIndex(graphLink.Links, id)
				if i != -1 {
					links = append(graphLink.Links[:i], graphLink.Links[i+1:]...)
				}

				_, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graphLink.Id}, bson.M{"$set": bson.M{"links": links}})
				if err != nil {
					return nil, err
				}
			}
		}

		if len(added) > 0 {
			var addedGraphs []models.GraphPoint
			addFilter := bson.M{
				"_id": bson.M{
					"$in": added,
				},
			}

			cur, err := graphsCol.Find(ctx, addFilter)
			if err != nil {
				return nil, err
			}

			err = cur.All(ctx, &addedGraphs)
			if err != nil {
				return nil, err
			}

			cur.Close(ctx)

			if len(addedGraphs) < len(added) {
				return nil, errors.New("some graph point is missing in database")
			}

			for _, graphLink := range addedGraphs {
				if graphLink.Institute != body.Institute {
					return nil, errors.New("one of added graph points in links has different institute value")
				}

				graphLink.Links = append(graphLink.Links, id)
				_, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graphLink.Id}, bson.M{"$set": graphLink})
				if err != nil {
					return nil, err
				}
			}
		}

		//Stair
		if body.StairId != graph.StairId {
			if graph.StairId != "" {
				var oldStair models.Stair
				err = stairsCol.FindOne(ctx, bson.M{"_id": graph.StairId}).Decode(&oldStair)
				if err != nil {
					return nil, err
				}

				stairLinks := oldStair.Links
				i := utils.GetIndex(oldStair.Links, id)
				if i != -1 {
					stairLinks = append(oldStair.Links[:i], oldStair.Links[i+1:]...)
				}

				_, err = stairsCol.UpdateOne(ctx, bson.M{"_id": graph.StairId}, bson.M{"$set": bson.M{"links": stairLinks}})
				if err != nil {
					return nil, err
				}
			}

			if body.StairId != "" {
				var newStair models.Stair
				err = stairsCol.FindOne(ctx, bson.M{"_id": body.StairId}).Decode(&newStair)
				if err != nil {
					return nil, err
				}

				if newStair.Institute != body.Institute {
					return nil, errors.New("graph point and stair have different institute value")
				} else {
					newStair.Links = append(newStair.Links, id)
					_, err = stairsCol.UpdateOne(ctx, bson.M{"_id": body.StairId}, bson.M{"$set": newStair})
					if err != nil {
						return nil, err
					}

					i := utils.GetIndex(body.Types, "stair")
					if i == -1 {
						body.Types = append(body.Types, "stair")
					}
				}
			} else {
				i := utils.GetIndex(body.Types, "stair")
				if i != -1 {
					body.Types = append(body.Types[:i], body.Types[i+1:]...)
				}
			}
		}

		//На всякий пожарный
		if body.StairId == "" && slices.Contains(body.Types, "stair") {
			return nil, errors.New("there is no stairId in body with type \"stair\"")
		}

		if body.StairId != "" && !slices.Contains(body.Types, "stair") {
			return nil, errors.New("there is stairId in body without type \"stair\"")
		}

		// newBody := utils.GraphPointPutToGraphPoint(body, id)

		_, err = graphsCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
		if err != nil {
			log.Println("1")
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	return err
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
		err := graphsCol.FindOne(ctx, graphFilter).Decode(&graph)
		if err != nil {
			return nil, err
		}

		var floor models.Floor
		if graph.Institute != "" {
			floorFilter := bson.M{
				"floor":     graph.Floor,
				"institute": graph.Institute,
			}
			err = floorsCol.FindOne(ctx, floorFilter).Decode(&floor)

			if err == nil {
				newGraphs := floor.Graph
				graphIndex := utils.GetIndex(floor.Graph, graph.Id)
				if graphIndex != -1 {
					newGraphs = append(floor.Graph[:graphIndex], floor.Graph[graphIndex+1:]...)
				}

				_, err = floorsCol.UpdateOne(ctx, bson.M{"_id": floor.Id}, bson.M{"$set": bson.M{"graph": newGraphs}})
				if err != nil {
					return nil, err
				}
			}
		}

		if graph.StairId != "" {
			stairFilter := bson.M{
				"stairPoint": graph.StairId,
			}
			var stair models.Stair
			err = stairsCol.FindOne(ctx, stairFilter).Decode(&stair)
			if err != nil {
				return nil, err
			}
			log.Println(stair.Id)

			newLinks := stair.Links
			linkIndex := utils.GetIndex(stair.Links, graph.Id)
			if linkIndex != -1 {
				newLinks = append(stair.Links[:linkIndex], stair.Links[linkIndex+1:]...)
			}
			_, err = stairsCol.UpdateOne(ctx, stairFilter, bson.M{"$set": bson.M{"links": newLinks}})
			if err != nil {
				return nil, err
			}
		}

		if len(graph.Links) > 0 {
			var linkGraphs []models.GraphPoint
			filter := bson.M{
				"_id": bson.M{
					"$in": graph.Links,
				},
			}

			cur, err := graphsCol.Find(ctx, filter)
			if err != nil {
				return nil, err
			}

			err = cur.All(ctx, &linkGraphs)
			if err != nil {
				return nil, err
			}

			cur.Close(ctx)

			for _, linkGraph := range linkGraphs {
				newLinks := linkGraph.Links
				linkIndex := utils.GetIndex(linkGraph.Links, id)
				if linkIndex != -1 {
					newLinks = append(linkGraph.Links[:linkIndex], linkGraph.Links[linkIndex+1:]...)
				}
				_, err = graphsCol.UpdateOne(ctx, bson.M{"_id": linkGraph.Id}, bson.M{"$set": bson.M{"links": newLinks}})
				if err != nil {
					return nil, err
				}
			}
		}

		_, err = graphsCol.DeleteOne(ctx, graphFilter)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	return err
}
