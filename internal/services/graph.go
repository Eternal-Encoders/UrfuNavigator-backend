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

func (s *Services) PostGraphs(context context.Context, graphs []*models.GraphPoint) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		newValue := make([]interface{}, len(graphs))

		for i := range graphs {
			newValue[i] = graphs[i]
		}

		if len(newValue) > 0 {
			res := s.Store.InsertManyGraphs(ctx, graphs)
			if res.Error != nil {
				return res, nil
			}
		}

		res := s.PostStairs(ctx, graphs)
		if res.Error != nil {
			log.Println("stair error")
			return res, nil
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}

func (s *Services) GetGraph(context context.Context, id string) (models.GraphPoint, models.ResponseType) {
	return s.Store.GetGraph(context, []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: id}})
}

func (s *Services) GetAllGraphs(context context.Context) ([]models.GraphPoint, models.ResponseType) {
	return s.Store.GetAllGraphs(context)
}

func (s *Services) UpdateGraph(context context.Context, body models.GraphPointPut, id string) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		// graphFilter := bson.M{
		// 	"_id": id,
		// }
		// var graph models.GraphPoint

		graphFilter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: id}}
		graph, res := s.Store.GetGraph(ctx, graphFilter)
		if res.Error != nil {
			return res, nil
		}

		// err := graphsCol.FindOne(ctx, graphFilter).Decode(&graph)
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + id)}, err
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}, err
		// 	}
		// }

		//Floors and institute

		var oldFloor models.Floor
		if graph.Institute != "" {
			oldFloorFilter := []models.Query{models.Query{ParamName: "floor", Type: "int", IntValue: graph.Floor},
				models.Query{ParamName: "institute", Type: "string", StringValue: graph.Institute}}
			oldFloor, res = s.Store.GetFloor(ctx, oldFloorFilter)
			if res.Error != nil {
				return res, nil
			}

			// oldFloorFilter := bson.M{
			// 	"floor":     graph.Floor,
			// 	"institute": graph.Institute,
			// }
			// err = floorsCol.FindOne(ctx, oldFloorFilter).Decode(&oldFloor)
			// if err != nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified institute and floor: " + graph.Institute + ", " + strconv.Itoa(graph.Floor))}, err
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}, err
			// 	}
			// }
		}

		var newFloor models.Floor
		if body.Institute != "" {
			newFloorFilter := []models.Query{models.Query{ParamName: "floor", Type: "int", IntValue: body.Floor},
				models.Query{ParamName: "institute", Type: "string", StringValue: body.Institute}}
			newFloor, res = s.Store.GetFloor(ctx, newFloorFilter)
			if res.Error != nil {
				return res, nil
			}

			// newFloorFilter := bson.M{
			// 	"floor":     body.Floor,
			// 	"institute": body.Institute,
			// }
			// err = floorsCol.FindOne(ctx, newFloorFilter).Decode(&newFloor)
			// if err != nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified institute and floor: " + graph.Institute + ", " + strconv.Itoa(graph.Floor))}, err
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}, err
			// 	}
			// }
		}

		if newFloor.Id != oldFloor.Id {
			if body.Institute != "" {
				newFloor.Graph = append(newFloor.Graph, id)

				filter := []models.Query{models.Query{ParamName: "_id", Type: "ObjectID", ObjectIDValue: newFloor.Id}}

				res = s.Store.UpdateFloor(ctx, filter, utils.FloorToFloorPut(newFloor))
				if res.Error != nil {
					return res, nil
				}

				// _, err = floorsCol.UpdateOne(ctx, bson.M{"_id": newFloor.Id}, bson.M{"$set": bson.M{"graph": newFloor.Graph}})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified id: " + string(newFloor.Id.Hex()))}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, nil
				// 	}
				// }
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

				utils.FloorToFloorPut(oldFloor)

				filter := []models.Query{models.Query{ParamName: "_id", Type: "ObjectID", ObjectIDValue: oldFloor.Id}}

				res = s.Store.UpdateFloor(ctx, filter, utils.FloorToFloorPut(oldFloor))
				if res.Error != nil {
					return res, nil
				}

				// _, err = floorsCol.UpdateOne(ctx, bson.M{"_id": oldFloor.Id}, bson.M{"$set": oldFloor})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified id: " + string(oldFloor.Id.Hex()))}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }
			}
		}

		added, deleted := utils.GetAddedDeleted(graph.Links, body.Links)

		if len(deleted) > 0 {
			var deletedGraphs []models.GraphPoint
			// delFilter := bson.M{
			// 	"_id": bson.M{
			// 		"$in": deleted,
			// 	},
			// }

			deletedGraphs, res = s.Store.GetManyGraphsByIds(ctx, deleted)
			if res.Error != nil {
				return res, nil
			}

			// cur, err := graphsCol.Find(ctx, delFilter)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// err = cur.All(ctx, &deletedGraphs)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// cur.Close(ctx)

			for _, graphLink := range deletedGraphs {
				links := graphLink.Links
				i := utils.GetIndex(graphLink.Links, id)
				if i != -1 {
					links = append(graphLink.Links[:i], graphLink.Links[i+1:]...)
				}

				graphLink.Links = links

				filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: graphLink.Id}}

				res = s.Store.UpdateGraph(ctx, filter, utils.GraphPointToGraphPointPut(graphLink))
				if res.Error != nil {
					return res, nil
				}

				// _, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graphLink.Id}, bson.M{"$set": bson.M{"links": links}})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + graphLink.Id)}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }
			}
		}

		if len(added) > 0 {
			var addedGraphs []models.GraphPoint

			addedGraphs, res = s.Store.GetManyGraphsByIds(ctx, added)
			if res.Error != nil {
				return res, nil
			}
			// addFilter := bson.M{
			// 	"_id": bson.M{
			// 		"$in": added,
			// 	},
			// }

			// cur, err := graphsCol.Find(ctx, addFilter)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// err = cur.All(ctx, &addedGraphs)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// cur.Close(ctx)

			if len(addedGraphs) < len(added) {
				return models.ResponseType{Type: 404, Error: errors.New("some graph point is missing in database")}, nil
			}

			for _, graphLink := range addedGraphs {
				if graphLink.Institute != body.Institute {
					return models.ResponseType{Type: 406, Error: errors.New("one of added graph points in links has different institute value")}, nil
				}

				graphLink.Links = append(graphLink.Links, id)

				filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: graphLink.Id}}

				res = s.Store.UpdateGraph(ctx, filter, utils.GraphPointToGraphPointPut(graphLink))
				if res.Error != nil {
					return res, nil
				}

				// _, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graphLink.Id}, bson.M{"$set": graphLink})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + graphLink.Id)}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }
			}
		}

		//Stair
		if body.StairId != graph.StairId {
			if graph.StairId != "" {
				var oldStair models.Stair

				filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: graph.StairId}}

				oldStair, res = s.Store.GetStair(ctx, filter)
				if res.Error != nil {
					return res, nil
				}

				// err = stairsCol.FindOne(ctx, bson.M{"_id": graph.StairId}).Decode(&oldStair)
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, nil
				// 	}
				// }

				stairLinks := oldStair.Links
				i := utils.GetIndex(oldStair.Links, id)
				if i != -1 {
					stairLinks = append(oldStair.Links[:i], oldStair.Links[i+1:]...)
				}

				oldStair.Links = stairLinks

				res = s.Store.UpdateStair(ctx, filter, oldStair)
				if res.Error != nil {
					return res, nil
				}

				// _, err = stairsCol.UpdateOne(ctx, bson.M{"_id": graph.StairId}, bson.M{"$set": bson.M{"links": stairLinks}})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }
			}

			if body.StairId != "" {
				var newStair models.Stair

				filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: body.StairId}}

				newStair, res = s.Store.GetStair(ctx, filter)
				if res.Error != nil {
					return res, nil
				}

				// err = stairsCol.FindOne(ctx, bson.M{"_id": body.StairId}).Decode(&newStair)
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + body.StairId)}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }

				if newStair.Institute != body.Institute {
					return models.ResponseType{Type: 406, Error: errors.New("graph point and stair have different institute value")}, nil
				} else {
					newStair.Links = append(newStair.Links, id)

					res = s.Store.UpdateStair(ctx, filter, newStair)
					if res.Error != nil {
						return res, nil
					}

					// _, err = stairsCol.UpdateOne(ctx, bson.M{"_id": body.StairId}, bson.M{"$set": newStair})
					// if err != nil {
					// 	if err == mongo.ErrNoDocuments {
					// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + body.StairId)}, err
					// 	} else {
					// 		return models.ResponseType{Type: 500, Error: err}, err
					// 	}
					// }

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
			return models.ResponseType{Type: 406, Error: errors.New("there is no stairId in body with type \"stair\"")}, nil
		}

		if body.StairId != "" && !slices.Contains(body.Types, "stair") {
			return models.ResponseType{Type: 406, Error: errors.New("there is stairId in body without type \"stair\"")}, nil
		}

		// newBody := utils.GraphPointPutToGraphPoint(body, id)

		filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: id}}

		res = s.Store.UpdateGraph(ctx, filter, body)
		if res.Error != nil {
			return res, nil
		}

		// _, err = graphsCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": body})
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + id)}, err
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}, err
		// 	}
		// }

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}

func (s *Services) DeleteGraph(context context.Context, id string) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		// graphFilter := bson.M{
		// 	"_id": id,
		// }
		var graph models.GraphPoint

		graphFilter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: id}}
		graph, res := s.Store.GetGraph(ctx, graphFilter)
		if res.Error != nil {
			return res, nil
		}

		// err := graphsCol.FindOne(ctx, graphFilter).Decode(&graph)
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + id)}, err
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}, err
		// 	}
		// }

		var floor models.Floor
		if graph.Institute != "" {
			// floorFilter := bson.M{
			// 	"floor":     graph.Floor,
			// 	"institute": graph.Institute,
			// }

			floorFilter := []models.Query{models.Query{ParamName: "floor", Type: "int", IntValue: graph.Floor},
				models.Query{ParamName: "institute", Type: "string", StringValue: graph.Institute}}
			floor, res = s.Store.GetFloor(ctx, floorFilter)

			// err = floorsCol.FindOne(ctx, floorFilter).Decode(&floor)

			if res.Error == nil {
				newGraphs := floor.Graph
				graphIndex := utils.GetIndex(floor.Graph, graph.Id)
				if graphIndex != -1 {
					newGraphs = append(floor.Graph[:graphIndex], floor.Graph[graphIndex+1:]...)
				}

				floor.Graph = newGraphs

				filter := []models.Query{models.Query{ParamName: "_id", Type: "ObjectID", ObjectIDValue: floor.Id}}
				res = s.Store.UpdateFloor(ctx, filter, utils.FloorToFloorPut(floor))
				if res.Error != nil {
					return res, nil
				}

				// _, err = floorsCol.UpdateOne(ctx, bson.M{"_id": floor.Id}, bson.M{"$set": bson.M{"graph": newGraphs}})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified id: " + floor.Id.Hex())}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }
			}
		}

		if graph.StairId != "" {
			// stairFilter := bson.M{
			// 	"stairPoint": graph.StairId,
			// }
			var stair models.Stair

			stairFilter := []models.Query{models.Query{ParamName: "stairPoint", Type: "string", StringValue: graph.StairId}}

			stair, res = s.Store.GetStair(ctx, stairFilter)
			if res.Error != nil {
				return res, nil
			}

			// err = stairsCol.FindOne(ctx, stairFilter).Decode(&stair)
			// if err != nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}, err
			// 	}
			// }

			newLinks := stair.Links
			linkIndex := utils.GetIndex(stair.Links, graph.Id)
			if linkIndex != -1 {
				newLinks = append(stair.Links[:linkIndex], stair.Links[linkIndex+1:]...)
			}

			stair.Links = newLinks

			res = s.Store.UpdateStair(ctx, stairFilter, stair)
			if res.Error != nil {
				return res, nil
			}

			// _, err = stairsCol.UpdateOne(ctx, stairFilter, bson.M{"$set": bson.M{"links": newLinks}})
			// if err != nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}, err
			// 	}
			// }
		}

		if len(graph.Links) > 0 {
			var linkGraphs []models.GraphPoint

			linkGraphs, res = s.Store.GetManyGraphsByIds(ctx, graph.Links)
			if res.Error != nil {
				return res, nil
			}

			// filter := bson.M{
			// 	"_id": bson.M{
			// 		"$in": graph.Links,
			// 	},
			// }

			// cur, err := graphsCol.Find(ctx, filter)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// err = cur.All(ctx, &linkGraphs)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// cur.Close(ctx)

			for _, linkGraph := range linkGraphs {
				newLinks := linkGraph.Links
				linkIndex := utils.GetIndex(linkGraph.Links, id)
				if linkIndex != -1 {
					newLinks = append(linkGraph.Links[:linkIndex], linkGraph.Links[linkIndex+1:]...)
				}

				linkGraph.Links = newLinks
				filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: linkGraph.Id}}

				res = s.Store.UpdateGraph(ctx, filter, utils.GraphPointToGraphPointPut(linkGraph))
				if res.Error != nil {
					return res, nil
				}
				// _, err = graphsCol.UpdateOne(ctx, bson.M{"_id": linkGraph.Id}, bson.M{"$set": bson.M{"links": newLinks}})
				// if err != nil {
				// 	if err == mongo.ErrNoDocuments {
				// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + linkGraph.Id)}, err
				// 	} else {
				// 		return models.ResponseType{Type: 500, Error: err}, err
				// 	}
				// }
			}
		}

		res = s.Store.DeleteGraph(ctx, graphFilter)
		if res.Error != nil {
			return res, nil
		}

		// _, err = graphsCol.DeleteOne(ctx, graphFilter)
		// if err != nil {
		// 	return models.ResponseType{Type: 500, Error: err}, err
		// }

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}
