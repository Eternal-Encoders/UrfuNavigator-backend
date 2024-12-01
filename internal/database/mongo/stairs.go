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

func (s *MongoDB) PostStairs(context context.Context, graphs []*models.GraphPoint) error {
	collection := s.Database.Collection("stairs")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return err
	}

	defer session.EndSession(context)

	_, err = session.WithTransaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		var stairsForInsert []models.Stair
		var stairIds []string
		var existingStairs []models.Stair

		for _, graph := range graphs {
			if graph.StairId != "" {
				if !slices.Contains(stairIds, graph.StairId) {
					stairIds = append(stairIds, graph.StairId)

					filter := bson.M{"_id": graph.StairId}
					var stair models.Stair

					if err := collection.FindOne(ctx, filter).Decode(&stair); err == nil {
						stair.Links = append(stair.Links, graph.Id)
						existingStairs = append(existingStairs, stair)
					} else {
						stair = models.Stair{
							Id:         graph.StairId,
							StairPoint: graph.StairId,
							Institute:  graph.Institute,
							Links:      []string{graph.Id},
						}

						stairsForInsert = append(stairsForInsert, stair)
					}
				} else {
					found := false

					for _, stair := range existingStairs {
						if stair.Id == graph.StairId {
							stair.Links = append(stair.Links, graph.Id)
							found = true
							break
						}
					}

					if !found {
						for _, stair := range stairsForInsert {
							if stair.Id == graph.StairId {
								stair.Links = append(stair.Links, graph.Id)
								found = true
								break
							}
						}
					}
				}
			}
		}

		insertRes := make([]interface{}, len(stairsForInsert))
		for i := range stairsForInsert {
			insertRes[i] = stairsForInsert[i]
		}

		if len(insertRes) > 0 {
			_, err := collection.InsertMany(ctx, insertRes)
			if err != nil {
				log.Println("stair insertion error")
				return nil, err
			}
		}

		for _, stair := range existingStairs {
			_, err := collection.UpdateOne(ctx, bson.M{"_id": stair.Id}, bson.M{"$set": stair})
			if err != nil {
				log.Println("stair update error")
				return nil, err
			}
		}

		return nil, nil
	}, txnOptions)

	return err
}

func (s *MongoDB) GetStair(id string) (models.Stair, error) {
	collection := s.Database.Collection("stairs")

	filter := bson.M{
		"_id": id,
	}

	var result models.Stair
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func (s *MongoDB) GetAllStairs() ([]models.Stair, error) {
	collection := s.Database.Collection("stairs")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var result []models.Stair
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return result, nil
}

func (s *MongoDB) UpdateStair(context context.Context, body models.Stair, id string) error {
	stairsCol := s.Database.Collection("stairs")
	graphsCol := s.Database.Collection("graph_points")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return err
	}

	defer session.EndSession(context)

	_, err = session.WithTransaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		stairFilter := bson.M{
			"_id": id,
		}
		var oldStair models.Stair
		err := stairsCol.FindOne(ctx, stairFilter).Decode(&oldStair)
		if err != nil {
			return nil, err
		}

		if body.Id != id || body.StairPoint != id {
			return nil, errors.New("the specified id does not match with body id")
		}

		// if body.StairPoint != body.Id {
		// 	if oldStair.StairPoint != body.StairPoint {
		// 		body.Id = body.StairPoint
		// 	} else if oldStair.Id != body.Id {
		// 		body.StairPoint = body.Id
		// 	}
		// }

		added, deleted := utils.GetAddedDeleted(oldStair.Links, body.Links)

		if body.Institute != oldStair.Institute {
			err = utils.UpdateGraphsStairs(ctx, body.Links, body, graphsCol, stairsCol)
			if err != nil {
				return nil, err
			}
		} else {
			err = utils.UpdateGraphsStairs(ctx, added, body, graphsCol, stairsCol)
			if err != nil {
				return nil, err
			}
		}

		if len(deleted) > 0 {
			var deletedGraphs []models.GraphPoint
			filter := bson.M{
				"_id": bson.M{
					"$in": deleted,
				},
			}

			cur, err := graphsCol.Find(ctx, filter)
			if err != nil {
				return nil, err
			}

			err = cur.All(ctx, &deletedGraphs)
			if err != nil {
				return nil, err
			}

			cur.Close(ctx)

			for _, graph := range deletedGraphs {
				typeIndex := utils.GetIndex(graph.Types, "stair")
				newTypes := graph.Types
				if typeIndex != -1 {
					newTypes = append(graph.Types[:typeIndex], graph.Types[typeIndex+1:]...)
				}
				_, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graph.Id}, bson.M{"$set": bson.M{"stairId": "", "types": newTypes}})
				if err != nil {
					return nil, err
				}
			}
		}

		_, err = stairsCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	return err
}

func (s *MongoDB) DeleteStair(id string) error {
	stairsCol := s.Database.Collection("stairs")
	graphsCol := s.Database.Collection("graph_points")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return err
	}

	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		stairFilter := bson.M{
			"_id": id,
		}
		var stair models.Stair
		err := stairsCol.FindOne(ctx, stairFilter).Decode(&stair)
		if err != nil {
			return nil, err
		}
		log.Println(stair.StairPoint)

		if len(stair.Links) > 0 {
			graphFilter := bson.M{
				"_id": bson.M{"$in": stair.Links},
			}
			cursor, err := graphsCol.Find(ctx, graphFilter)
			if err != nil {
				return nil, err
			}

			defer cursor.Close(ctx)

			var graphs []models.GraphPoint
			decodeErr := cursor.All(ctx, &graphs)
			if decodeErr != nil {
				return nil, decodeErr
			}
			// log.Println(graphs)

			for _, v := range graphs {
				log.Println(v)
				typeIndex := utils.GetIndex(v.Types, "stair")
				newTypes := v.Types
				if typeIndex != -1 {
					newTypes = append(v.Types[:typeIndex], v.Types[typeIndex+1:]...)
				}

				_, err = graphsCol.UpdateOne(ctx, bson.M{"_id": v.Id}, bson.M{"$set": bson.M{"stairId": "", "types": newTypes}})
				if err != nil {
					return nil, err
				}
			}
		}

		_, err = stairsCol.DeleteOne(ctx, stairFilter)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	return err
}
