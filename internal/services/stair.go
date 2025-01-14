package services

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"errors"
	"log"
	"slices"

	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Services) PostStairs(context context.Context, graphs []*models.GraphPoint) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		var stairsForInsert []*models.Stair
		var stairIds []string
		var existingStairs []models.Stair

		for _, graph := range graphs {
			if graph.StairId != "" {
				if !slices.Contains(stairIds, graph.StairId) {
					stairIds = append(stairIds, graph.StairId)

					// filter := bson.M{"_id": graph.StairId}
					var stair models.Stair

					// if err := collection.FindOne(ctx, filter).Decode(&stair); err == nil {
					// 	stair.Links = append(stair.Links, graph.Id)
					// 	existingStairs = append(existingStairs, stair)
					// }

					filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: graph.StairId}}
					stair, res := s.Store.GetStair(ctx, filter)
					if res.Error == nil {
						stair.Links = append(stair.Links, graph.Id)
						existingStairs = append(existingStairs, stair)
					} else {
						stair = models.Stair{
							Id:         graph.StairId,
							StairPoint: graph.StairId,
							Institute:  graph.Institute,
							Links:      []string{graph.Id},
						}

						stairsForInsert = append(stairsForInsert, &stair)
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
			res := s.Store.InsertManyStairs(ctx, stairsForInsert)
			if res.Error != nil {
				return res, nil
			}
			// _, err := collection.InsertMany(ctx, insertRes)
			// if err != nil {
			// 	log.Println("stair insertion error")
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }
		}

		for _, stair := range existingStairs {
			filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: stair.Id}}
			res := s.Store.UpdateStair(ctx, filter, stair)
			if res.Error != nil {
				return res, nil
			}
		}

		// for _, stair := range existingStairs {
		// 	_, err := collection.UpdateOne(ctx, bson.M{"_id": stair.Id}, bson.M{"$set": stair})
		// 	if err != nil {
		// 		log.Println("stair update error")
		// 		if err == mongo.ErrNoDocuments {
		// 			return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + stair.Id)}, err
		// 		} else {
		// 			return models.ResponseType{Type: 500, Error: err}, err
		// 		}
		// 	}
		// }

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}

func (s *Services) UpdateStair(context context.Context, body models.Stair, id string) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		// stairFilter := bson.M{
		// 	"_id": id,
		// }
		var oldStair models.Stair
		stairFilter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: id}}

		oldStair, res := s.Store.GetStair(ctx, stairFilter)
		if res.Error != nil {
			return res, nil
		}

		// err := stairsCol.FindOne(ctx, stairFilter).Decode(&oldStair)
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + id)}, err
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}, err
		// 	}
		// }

		if body.Id != id || body.StairPoint != id {
			return models.ResponseType{Type: 406, Error: errors.New("the specified id does not match with body id")}, nil
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
			res := utils.UpdateGraphsStairs(ctx, body.Links, body, s.Store)
			if res.Error != nil {
				return res, res.Error
			}
		} else {
			res := utils.UpdateGraphsStairs(ctx, added, body, s.Store)
			if res.Error != nil {
				return res, res.Error
			}
		}

		if len(deleted) > 0 {
			var deletedGraphs []models.GraphPoint

			deletedGraphs, res := s.Store.GetManyGraphsByIds(ctx, deleted)
			if res.Error != nil {
				return res, res.Error
			}
			// filter := bson.M{
			// 	"_id": bson.M{
			// 		"$in": deleted,
			// 	},
			// }

			// cur, err := graphsCol.Find(ctx, filter)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// err = cur.All(ctx, &deletedGraphs)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// cur.Close(ctx)

			for _, graph := range deletedGraphs {
				typeIndex := utils.GetIndex(graph.Types, "stair")
				newTypes := graph.Types
				if typeIndex != -1 {
					newTypes = append(graph.Types[:typeIndex], graph.Types[typeIndex+1:]...)
				}

				graph.Types = newTypes
				graph.StairId = ""
				filter := []models.Query{{ParamName: "_id", Type: "string", StringValue: graph.Id}}

				res = s.Store.UpdateGraph(ctx, filter, utils.GraphPointToGraphPointPut(graph))
				if res.Error != nil {
					return res, res.Error
				}

				// _, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graph.Id}, bson.M{"$set": bson.M{"stairId": "", "types": newTypes}})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + graph.Id)}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }
			}
		}

		filter := []models.Query{{ParamName: "_id", Type: "string", StringValue: id}}
		res = s.Store.UpdateStair(ctx, filter, body)
		if res.Error != nil {
			return res, res.Error
		}

		// _, err = stairsCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + id)}, err
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}, err
		// 	}
		// }

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}

func (s *Services) DeleteStair(context context.Context, id string) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		// stairFilter := bson.M{
		// 	"_id": id,
		// }
		var stair models.Stair
		stairFilter := []models.Query{{ParamName: "_id", Type: "string", StringValue: id}}
		stair, res := s.Store.GetStair(ctx, stairFilter)
		if res.Error != nil {
			return res, res.Error
		}

		// err := stairsCol.FindOne(ctx, stairFilter).Decode(&stair)
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + id)}, err
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}, err
		// 	}
		// }

		if len(stair.Links) > 0 {
			graphs, res := s.Store.GetManyGraphsByIds(ctx, stair.Links)
			if res.Error != nil {
				return res, res.Error
			}

			// graphFilter := bson.M{
			// 	"_id": bson.M{"$in": stair.Links},
			// }
			// cursor, err := graphsCol.Find(ctx, graphFilter)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// defer cursor.Close(ctx)

			// var graphs []models.GraphPoint
			// decodeErr := cursor.All(ctx, &graphs)
			// if decodeErr != nil {
			// 	return models.ResponseType{Type: 500, Error: decodeErr}, decodeErr
			// }

			for _, v := range graphs {
				log.Println(v)
				typeIndex := utils.GetIndex(v.Types, "stair")
				newTypes := v.Types
				if typeIndex != -1 {
					newTypes = append(v.Types[:typeIndex], v.Types[typeIndex+1:]...)
				}

				v.Types = newTypes
				v.StairId = ""
				filter := []models.Query{{ParamName: "_id", Type: "string", StringValue: v.Id}}

				res = s.Store.UpdateGraph(ctx, filter, utils.GraphPointToGraphPointPut(v))
				if res.Error != nil {
					return res, res.Error
				}

				// _, err = graphsCol.UpdateOne(ctx, bson.M{"_id": v.Id}, bson.M{"$set": bson.M{"stairId": "", "types": newTypes}})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + v.Id)}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }
			}
		}

		res = s.Store.DeleteGraph(ctx, stairFilter)
		if res.Error != nil {
			return res, res.Error
		}

		// _, err = stairsCol.DeleteOne(ctx, stairFilter)
		// if err != nil {
		// 	return models.ResponseType{Type: 500, Error: err}, err
		// }

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}

func (s *Services) GetStair(context context.Context, id string) (models.Stair, models.ResponseType) {
	return s.Store.GetStair(context, []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: id}})
}

func (s *Services) GetAllStairs(context context.Context) ([]models.Stair, models.ResponseType) {
	return s.Store.GetAllStairs(context)
}
