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

func (s *MongoDB) PostStairs(context context.Context, graphs []*models.GraphPoint) models.ResponseType {
	collection := s.Database.Collection("stairs")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return models.ResponseType{Type: 500, Error: err}
	}

	defer session.EndSession(context)

	res, _ := session.WithTransaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
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
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}

		for _, stair := range existingStairs {
			_, err := collection.UpdateOne(ctx, bson.M{"_id": stair.Id}, bson.M{"$set": stair})
			if err != nil {
				log.Println("stair update error")
				if err == mongo.ErrNoDocuments {
					return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + stair.Id)}, err
				} else {
					return models.ResponseType{Type: 500, Error: err}, err
				}
			}
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	}, txnOptions)

	result := res.(models.ResponseType)
	return models.ResponseType{Type: result.Type, Error: result.Error}
}

func (s *MongoDB) GetStair(id string) (models.Stair, models.ResponseType) {
	collection := s.Database.Collection("stairs")

	filter := bson.M{
		"_id": id,
	}

	var result models.Stair
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Stair{}, models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + id)}
		} else {
			return models.Stair{}, models.ResponseType{Type: 500, Error: err}
		}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetAllStairs() ([]models.Stair, models.ResponseType) {
	collection := s.Database.Collection("stairs")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context.TODO())

	var result []models.Stair
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: decodeErr}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) UpdateStair(context context.Context, body models.Stair, id string) models.ResponseType {
	stairsCol := s.Database.Collection("stairs")
	graphsCol := s.Database.Collection("graph_points")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return models.ResponseType{Type: 500, Error: err}
	}

	defer session.EndSession(context)

	res, _ := session.WithTransaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		stairFilter := bson.M{
			"_id": id,
		}
		var oldStair models.Stair
		err := stairsCol.FindOne(ctx, stairFilter).Decode(&oldStair)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + id)}, err
			} else {
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}

		if body.Id != id || body.StairPoint != id {
			err = errors.New("the specified id does not match with body id")
			return models.ResponseType{Type: 406, Error: err}, err
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
			res := utils.UpdateGraphsStairs(ctx, body.Links, body, graphsCol, stairsCol)
			if res.Error != nil {
				return res, res.Error
			}
		} else {
			res := utils.UpdateGraphsStairs(ctx, added, body, graphsCol, stairsCol)
			if res.Error != nil {
				return res, res.Error
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
				return models.ResponseType{Type: 500, Error: err}, err
			}

			err = cur.All(ctx, &deletedGraphs)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
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
					if err == mongo.ErrNoDocuments {
						return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + graph.Id)}, err
					} else {
						return models.ResponseType{Type: 500, Error: err}, err
					}
				}
			}
		}

		_, err = stairsCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + id)}, err
			} else {
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	}, txnOptions)

	result := res.(models.ResponseType)
	return models.ResponseType{Type: result.Type, Error: result.Error}
}

func (s *MongoDB) DeleteStair(id string) models.ResponseType {
	stairsCol := s.Database.Collection("stairs")
	graphsCol := s.Database.Collection("graph_points")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return models.ResponseType{Type: 500, Error: err}
	}

	defer session.EndSession(context.TODO())

	res, _ := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		stairFilter := bson.M{
			"_id": id,
		}
		var stair models.Stair
		err := stairsCol.FindOne(ctx, stairFilter).Decode(&stair)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + id)}, err
			} else {
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}
		log.Println(stair.StairPoint)

		if len(stair.Links) > 0 {
			graphFilter := bson.M{
				"_id": bson.M{"$in": stair.Links},
			}
			cursor, err := graphsCol.Find(ctx, graphFilter)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}

			defer cursor.Close(ctx)

			var graphs []models.GraphPoint
			decodeErr := cursor.All(ctx, &graphs)
			if decodeErr != nil {
				return models.ResponseType{Type: 500, Error: decodeErr}, decodeErr
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
					if err == mongo.ErrNoDocuments {
						return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + v.Id)}, err
					} else {
						return models.ResponseType{Type: 500, Error: err}, err
					}
				}
			}
		}

		_, err = stairsCol.DeleteOne(ctx, stairFilter)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, err
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	}, txnOptions)

	result := res.(models.ResponseType)
	return models.ResponseType{Type: result.Type, Error: result.Error}
}
