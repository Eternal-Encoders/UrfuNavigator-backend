package services

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"slices"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Services) PostFloor(context context.Context, file *multipart.FileHeader) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		var floorFromFile models.FloorFromFile

		f, err := file.Open()
		if err != nil {
			return models.ResponseType{Type: 400, Error: errors.New("Something went wrong while reading file")}, nil
		}
		defer f.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, f); err != nil {
			return models.ResponseType{Type: 500, Error: errors.New("Something went wrong while reading file into []bytes")}, nil
		}

		json.Unmarshal([]byte(buf.Bytes()), &floorFromFile)

		audArr := []*models.Auditorium{}
		for _, v := range floorFromFile.Audiences {
			audArr = append(audArr, v)
		}

		graphArr := []*models.GraphPoint{}
		graphKeysArr := []string{}
		for k, v := range floorFromFile.Graph {
			graphArr = append(graphArr, v)
			graphKeysArr = append(graphKeysArr, k)
		}

		floor := models.FloorRequest{
			Institute: floorFromFile.Institute,
			Floor:     floorFromFile.Floor,
			Width:     floorFromFile.Width,
			Height:    floorFromFile.Height,
			Service:   floorFromFile.Service,
			Audiences: audArr,
			Graph:     graphKeysArr,
		}

		filter := []models.Query{models.Query{ParamName: "name", Type: "string", StringValue: floor.Institute}}
		if _, res := s.Store.GetInstitute(ctx, filter); res.Error != nil {
			err = errors.New("there is no institute with specified name: " + floor.Institute)
			return models.ResponseType{Type: 404, Error: err}, nil
		}

		filter = []models.Query{models.Query{ParamName: "institute", Type: "string", StringValue: floor.Institute},
			models.Query{ParamName: "floor", Type: "int", IntValue: floor.Floor}}
		if _, res := s.Store.GetFloor(ctx, filter); res.Error == nil {
			err = errors.New("floor already exists")
			return models.ResponseType{Type: 406, Error: err}, nil
		}

		res := s.Store.InsertFloor(ctx, floor)
		if res.Error != nil {
			return res, nil
		}

		///

		res = s.PostGraphs(ctx, graphArr)
		if res.Error != nil {
			return res, nil
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}

func (s *Services) GetFloor(context context.Context, id string) (models.FloorResponse, models.ResponseType) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.FloorResponse{}, models.ResponseType{Type: 500, Error: err}
	}

	filter := []models.Query{models.Query{ParamName: "_id", Type: "ObjectID", ObjectIDValue: objId}}
	floorData, res := s.Store.GetFloor(context, filter)
	if res.Error != nil {
		log.Println(res)
		return models.FloorResponse{}, res
	}

	response := models.FloorResponse{
		Id:        floorData.Id.Hex(),
		Institute: floorData.Institute,
		Floor:     floorData.Floor,
		Width:     floorData.Width,
		Height:    floorData.Height,
		Audiences: floorData.Audiences,
		Service:   floorData.Service,
		Graph:     floorData.Graph,
	}

	return response, models.ResponseType{Type: 200, Error: nil}
}

func (s *Services) GetAllFloors(context context.Context) ([]models.FloorResponse, models.ResponseType) {
	floorData, res := s.Store.GetAllFloors(context)
	if res.Error != nil {
		log.Println(res)
		return []models.FloorResponse{}, res
	}

	response := []models.FloorResponse{}
	for _, floor := range floorData {
		response = append(response, models.FloorResponse{
			Id:        floor.Id.Hex(),
			Institute: floor.Institute,
			Floor:     floor.Floor,
			Width:     floor.Width,
			Height:    floor.Height,
			Audiences: floor.Audiences,
			Service:   floor.Service,
			Graph:     floor.Graph,
		})
	}

	return response, models.ResponseType{Type: 200, Error: nil}
}

func (s *Services) UpdateFloor(context context.Context, body models.FloorPut, id string) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {

		// objId, err := primitive.ObjectIDFromHex(id)
		// if err != nil {
		// 	return models.ResponseType{Type: 500, Error: err}, err
		// }

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, nil
		}

		filter := []models.Query{models.Query{ParamName: "_id", Type: "ObjectID", ObjectIDValue: objId}}
		// bson.M{
		// 	"_id": objId,
		// }

		var oldFloor models.Floor
		oldFloor, res := s.Store.GetFloor(ctx, filter)
		if res.Error != nil {
			return res, nil
		}
		// err = floorsCol.FindOne(ctx, filter).Decode(&oldFloor)
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified id: " + id)}, err
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}, err
		// 	}
		// }

		added, deleted := utils.GetAddedDeleted(oldFloor.Graph, body.Graph)
		var remained []string
		for _, v := range oldFloor.Graph {
			if slices.Contains(body.Graph, v) {
				remained = append(remained, v)
			}
		}

		if len(deleted) > 0 {
			var delGraphs []models.GraphPoint
			// delFilter := bson.M{
			// 	"_id": bson.M{
			// 		"$in": deleted,
			// 	},
			// }

			delGraphs, res := s.Store.GetManyGraphsByIds(ctx, deleted)
			if res.Error != nil {
				return res, nil
			}

			// cur, err := graphsCol.Find(ctx, delFilter)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// err = cur.All(ctx, &delGraphs)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// cur.Close(ctx)

			for _, delGraph := range delGraphs {
				putDelGraph := models.GraphPointPut{
					X:           delGraph.X,
					Y:           delGraph.Y,
					Links:       nil,
					Types:       nil,
					Names:       nil,
					Floor:       0,
					Institute:   "",
					Time:        delGraph.Time,
					Description: delGraph.Description,
					Info:        delGraph.Info,
					MenuId:      delGraph.MenuId,
					IsPassFree:  delGraph.IsPassFree,
					StairId:     "",
				}

				// filter = []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: delGraph.Id}}

				res = s.UpdateGraph(ctx, putDelGraph, delGraph.Id)
				if res.Error != nil {
					return res, nil
				}

				// res := s.UpdateGraph(ctx, putDelGraph, delGraph.Id)
				// if res.Error != nil {
				// 	log.Println("UpdateGraph error")
				// 	return res, err
				// }
			}
		}

		if oldFloor.Institute != body.Institute {
			var newInstitute models.Institute
			// instituteFilter := bson.M{
			// 	"name": body.Institute,
			// }
			instituteFilter := []models.Query{models.Query{ParamName: "name", Type: "string", StringValue: body.Institute}}

			newInstitute, res = s.Store.GetInstitute(ctx, instituteFilter)
			if res.Error != nil {
				return res, nil
			}

			// err = institutesCol.FindOne(ctx, instituteFilter).Decode(&newInstitute)
			// if err != nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified name: " + body.Institute)}, err
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}, err
			// 	}
			// }

			if newInstitute.MinFloor > body.Floor || newInstitute.MaxFloor < body.Floor {
				err = errors.New("floor is out of istitute floor bounds")
				return models.ResponseType{Type: 406, Error: err}, err
			}

			filter = []models.Query{models.Query{ParamName: "institute", Type: "string", StringValue: body.Institute},
				models.Query{ParamName: "floor", Type: "int", IntValue: body.Floor}}

			_, res = s.Store.GetFloor(ctx, filter)
			if res.Error != nil {
				return res, nil
			}

			// if err = floorsCol.FindOne(ctx, bson.M{"institute": body.Institute, "floor": body.Floor}).Err(); err == nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified institute and floor: " + body.Institute + ", " + strconv.Itoa(body.Floor))}, err
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}, err
			// 	}
			// }
		} else {
			var oldInstitute models.Institute

			filter = []models.Query{models.Query{ParamName: "name", Type: "string", StringValue: body.Institute}}
			oldInstitute, res = s.Store.GetInstitute(ctx, filter)
			if res.Error != nil {
				return res, nil
			}

			// err = institutesCol.FindOne(ctx, bson.M{"name": body.Institute}).Decode(&oldInstitute)
			// if err != nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified name: " + body.Institute)}, err
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}, err
			// 	}
			// }

			if oldInstitute.MinFloor > body.Floor || oldInstitute.MaxFloor < body.Floor {
				err = errors.New("floor is out of istitute floor bounds")
				return models.ResponseType{Type: 406, Error: err}, err
			}
		}

		if oldFloor.Institute != body.Institute || oldFloor.Floor != body.Floor {
			if len(remained) > 0 {
				var remainedGraphs []models.GraphPoint

				remainedGraphs, res = s.Store.GetManyGraphsByIds(ctx, remained)
				if res.Error != nil {
					return res, nil
				}

				// remFilter := bson.M{
				// 	"_id": bson.M{
				// 		"$in": remained,
				// 	},
				// }

				// cur, err := graphsCol.Find(ctx, remFilter)
				// if err != nil {
				// 	return models.ResponseType{Type: 500, Error: err}, err
				// }

				// err = cur.All(ctx, &remainedGraphs)
				// if err != nil {
				// 	return models.ResponseType{Type: 500, Error: err}, err
				// }

				// cur.Close(ctx)

				if len(remainedGraphs) < len(remained) {
					err = errors.New("some graph point is missing in database")
					return models.ResponseType{Type: 404, Error: err}, err
				}

				for _, graph := range remainedGraphs {
					graph.Floor = body.Floor
					graph.Institute = body.Institute

					if oldFloor.Institute != body.Institute {
						graph.StairId = ""

						i := utils.GetIndex(graph.Types, "stair")
						if i != -1 {
							graph.Types = append(graph.Types[:i], graph.Types[i+1:]...)
						}

						var oldStair models.Stair

						filter = []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: graph.StairId}}
						oldStair, res = s.Store.GetStair(ctx, filter)
						if res.Error != nil {
							return res, nil
						}
						// err = stairsCol.FindOne(ctx, bson.M{"_id": graph.StairId}).Decode(&oldStair)
						// if err != nil {
						// 	return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
						// }

						stairLinks := oldStair.Links
						i = utils.GetIndex(oldStair.Links, graph.Id)
						if i != -1 {
							stairLinks = append(oldStair.Links[:i], oldStair.Links[i+1:]...)
						}

						oldStair.Links = stairLinks

						filter = []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: graph.StairId}}
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

					filter = []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: graph.Id}}
					graphPut := models.GraphPointPut{
						X:           graph.X,
						Y:           graph.Y,
						Names:       graph.Names,
						Floor:       body.Floor,
						Institute:   body.Institute,
						Time:        graph.Time,
						Description: graph.Description,
						Info:        graph.Info,
						MenuId:      graph.MenuId,
						IsPassFree:  graph.IsPassFree,
						StairId:     "",
					}
					res = s.Store.UpdateGraph(ctx, filter, graphPut)
					if res.Error != nil {
						return res, nil
					}

					// _, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graph.Id}, bson.M{"$set": graph})
					// if err != nil {
					// 	if err == mongo.ErrNoDocuments {
					// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + graph.Id)}, err
					// 	} else {
					// 		return models.ResponseType{Type: 500, Error: err}, err
					// 	}
					// }
				}
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
				err = errors.New("some graph point is missing in database")
				return models.ResponseType{Type: 404, Error: err}, err
			}

			for _, graph := range addedGraphs {
				var putAddGraph models.GraphPointPut

				putAddGraph.Links = nil
				putAddGraph.StairId = ""
				i := utils.GetIndex(graph.Types, "stair")
				if i != -1 {
					putAddGraph.Types = append(graph.Types[:i], graph.Types[i+1:]...)
				}
				putAddGraph.X = graph.X
				putAddGraph.Y = graph.Y
				putAddGraph.Names = graph.Names
				putAddGraph.Floor = body.Floor
				putAddGraph.Institute = body.Institute
				putAddGraph.Time = graph.Time
				putAddGraph.Description = graph.Description
				putAddGraph.Info = graph.Info
				putAddGraph.MenuId = graph.MenuId
				putAddGraph.IsPassFree = graph.IsPassFree
				putAddGraph.StairId = ""

				res := s.UpdateGraph(ctx, putAddGraph, graph.Id)
				if res.Error != nil {
					return res, res.Error
				}
			}
		}

		for _, el := range body.Audiences {
			if !slices.Contains(body.Graph, el.PointId) && el.PointId != "" {
				err = errors.New("wrong PointId of one of the auditories")
				return models.ResponseType{Type: 404, Error: err}, err
			}
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}

func (s *Services) DeleteFloor(context context.Context, id string) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, err
		}
		// floorFilter := bson.M{
		// 	"_id": objId,
		// }

		floorFilter := []models.Query{models.Query{ParamName: "_id", Type: "ObjectID", ObjectIDValue: objId}}
		floor, res := s.Store.GetFloor(ctx, floorFilter)
		if res.Error != nil {
			return res, res.Error
		}

		// var floor models.Floor
		// err = floorsCol.FindOne(ctx, floorFilter).Decode(&floor)
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified id: " + id)}, err
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}, err
		// 	}
		// }

		if len(floor.Graph) > 0 {
			// var graphs []models.GraphPoint
			// filter := bson.M{
			// 	"_id": bson.M{
			// 		"$in": floor.Graph,
			// 	},
			// }

			graphs, res := s.Store.GetManyGraphsByIds(ctx, floor.Graph)
			if res.Error != nil {
				return res, res.Error
			}

			// cur, err := graphsCol.Find(ctx, filter)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// err = cur.All(ctx, &graphs)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }

			// cur.Close(ctx)

			for _, graph := range graphs {
				if graph.StairId != "" {
					stairFilter := []models.Query{models.Query{ParamName: "stairPoint", Type: "string", StringValue: graph.StairId}}

					stair, res := s.Store.GetStair(ctx, stairFilter)
					if res.Error != nil {
						return res, res.Error
					}

					// stairFilter := bson.M{
					// 	"stairPoint": graph.StairId,
					// }
					// var stair models.Stair
					// err = stairsCol.FindOne(ctx, stairFilter).Decode(&stair)
					// if err != nil {
					// 	if err == mongo.ErrNoDocuments {
					// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
					// 	} else {
					// 		return models.ResponseType{Type: 500, Error: err}, err
					// 	}
					// }

					linkIndex := utils.GetIndex(stair.Links, graph.Id)
					newLinks := stair.Links
					if linkIndex != -1 {
						newLinks = append(stair.Links[:linkIndex], stair.Links[linkIndex+1:]...)
					}

					stair.Links = newLinks
					res = s.Store.UpdateStair(ctx, stairFilter, stair)
					if res.Error != nil {
						return res, res.Error
					}

					// _, err = stairsCol.UpdateOne(ctx, stairFilter, bson.M{"$set": bson.M{"links": newLinks}})
					// if err != nil {
					// 	return models.ResponseType{Type: 500, Error: err}, err
					// }
				}
			}

			res = s.Store.DeleteManyGraphsByIds(ctx, floor.Graph)
			if res.Error != nil {
				return res, res.Error
			}

			// _, err = graphsCol.DeleteMany(ctx, filter)
			// if err != nil {
			// 	return models.ResponseType{Type: 500, Error: err}, err
			// }
		}

		res = s.Store.DeleteFloor(ctx, floorFilter)
		if res.Error != nil {
			return res, res.Error
		}

		// _, err = floorsCol.DeleteOne(ctx, floorFilter)
		// if err != nil {
		// 	return models.ResponseType{Type: 500, Error: err}, err
		// }

		return models.ResponseType{Type: 200, Error: nil}, nil
	})
}
