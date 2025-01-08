package database

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"errors"
	"log"
	"slices"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *MongoDB) PostFloor(floor models.FloorRequest, graphs []*models.GraphPoint) models.ResponseType {
	floorsCol := s.Database.Collection("floors")
	institutesCol := s.Database.Collection("institutes")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return models.ResponseType{Type: 500, Error: err}
	}

	defer session.EndSession(context.TODO())

	res, _ := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		if err := institutesCol.FindOne(ctx, bson.M{"name": floor.Institute}).Err(); err != nil {
			err = errors.New("there is no institute with specified name: " + floor.Institute)
			return models.ResponseType{Type: 404, Error: err}, err
		}

		if err := floorsCol.FindOne(ctx, bson.M{"institute": floor.Institute, "floor": floor.Floor}).Err(); err == nil {
			err = errors.New("floor already exists")
			return models.ResponseType{Type: 406, Error: err}, err
		}

		_, err := floorsCol.InsertOne(ctx, floor)
		if err != nil {
			log.Println("floor insertion error")
			return models.ResponseType{Type: 500, Error: err}, err
		}

		res := s.PostGraphs(ctx, graphs)
		if res.Error != nil {
			return res, res.Error
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	}, txnOptions)

	result := res.(models.ResponseType)
	return models.ResponseType{Type: result.Type, Error: result.Error}
}

func (s *MongoDB) GetFloor(id string) (models.Floor, models.ResponseType) {
	collection := s.Database.Collection("floors")

	log.Println(id)
	objId, err := primitive.ObjectIDFromHex(id)
	log.Println(objId)
	if err != nil {
		return models.Floor{}, models.ResponseType{Type: 500, Error: err}
	}

	filter := bson.M{
		"_id": objId,
	}

	var result models.Floor
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Floor{}, models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified id: " + id)}
		} else {
			return models.Floor{}, models.ResponseType{Type: 500, Error: err}
		}
	}
	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetAllFloors() ([]models.Floor, models.ResponseType) {
	collection := s.Database.Collection("floors")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context.TODO())

	var result []models.Floor
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: decodeErr}
	}

	// if len(result) == 0 {
	// 	return nil, models.ResponseType{Type: 404, Error: errors.New("there are no floors")}
	// }

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) UpdateFloor(body models.FloorPut, id string) models.ResponseType {
	graphsCol := s.Database.Collection("graph_points")
	floorsCol := s.Database.Collection("floors")
	institutesCol := s.Database.Collection("institutes")
	stairsCol := s.Database.Collection("stairs")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return models.ResponseType{Type: 500, Error: err}
	}

	defer session.EndSession(context.TODO())

	res, _ := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		objId, err := primitive.ObjectIDFromHex(id)
		log.Println(objId)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, err
		}

		filter := bson.M{
			"_id": objId,
		}

		var oldFloor models.Floor
		err = floorsCol.FindOne(ctx, filter).Decode(&oldFloor)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified id: " + id)}, err
			} else {
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}

		added, deleted := utils.GetAddedDeleted(oldFloor.Graph, body.Graph)
		var remained []string
		for _, v := range oldFloor.Graph {
			if slices.Contains(body.Graph, v) {
				remained = append(remained, v)
			}
		}

		if len(deleted) > 0 {
			var delGraphs []models.GraphPoint
			delFilter := bson.M{
				"_id": bson.M{
					"$in": deleted,
				},
			}

			cur, err := graphsCol.Find(ctx, delFilter)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}

			err = cur.All(ctx, &delGraphs)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}

			cur.Close(ctx)

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

				res := s.UpdateGraph(ctx, putDelGraph, delGraph.Id)
				if res.Error != nil {
					log.Println("UpdateGraph error")
					return res, err
				}
			}
		}

		if oldFloor.Institute != body.Institute {
			var newInstitute models.Institute
			instituteFilter := bson.M{
				"name": body.Institute,
			}

			err = institutesCol.FindOne(ctx, instituteFilter).Decode(&newInstitute)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified name: " + body.Institute)}, err
				} else {
					return models.ResponseType{Type: 500, Error: err}, err
				}
			}

			if newInstitute.MinFloor > body.Floor || newInstitute.MaxFloor < body.Floor {
				err = errors.New("floor is out of istitute floor bounds")
				return models.ResponseType{Type: 406, Error: err}, err
			}

			if err = floorsCol.FindOne(ctx, bson.M{"institute": body.Institute, "floor": body.Floor}).Err(); err == nil {
				if err == mongo.ErrNoDocuments {
					return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified institute and floor: " + body.Institute + ", " + strconv.Itoa(body.Floor))}, err
				} else {
					return models.ResponseType{Type: 500, Error: err}, err
				}
			}
		} else {
			var oldInstitute models.Institute

			err = institutesCol.FindOne(ctx, bson.M{"name": body.Institute}).Decode(&oldInstitute)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified name: " + body.Institute)}, err
				} else {
					return models.ResponseType{Type: 500, Error: err}, err
				}
			}

			if oldInstitute.MinFloor > body.Floor || oldInstitute.MaxFloor < body.Floor {
				err = errors.New("floor is out of istitute floor bounds")
				return models.ResponseType{Type: 406, Error: err}, err
			}
		}

		if oldFloor.Institute != body.Institute || oldFloor.Floor != body.Floor {
			if len(remained) > 0 {
				var remainedGraphs []models.GraphPoint
				remFilter := bson.M{
					"_id": bson.M{
						"$in": remained,
					},
				}

				cur, err := graphsCol.Find(ctx, remFilter)
				if err != nil {
					return models.ResponseType{Type: 500, Error: err}, err
				}

				err = cur.All(ctx, &remainedGraphs)
				if err != nil {
					return models.ResponseType{Type: 500, Error: err}, err
				}

				cur.Close(ctx)

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
						err = stairsCol.FindOne(ctx, bson.M{"_id": graph.StairId}).Decode(&oldStair)
						if err != nil {
							return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
						}

						stairLinks := oldStair.Links
						i = utils.GetIndex(oldStair.Links, graph.Id)
						if i != -1 {
							stairLinks = append(oldStair.Links[:i], oldStair.Links[i+1:]...)
						}

						_, err = stairsCol.UpdateOne(ctx, bson.M{"_id": graph.StairId}, bson.M{"$set": bson.M{"links": stairLinks}})
						if err != nil {
							if err == mongo.ErrNoDocuments {
								return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
							} else {
								return models.ResponseType{Type: 500, Error: err}, err
							}
						}
					}

					_, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graph.Id}, bson.M{"$set": graph})
					if err != nil {
						if err == mongo.ErrNoDocuments {
							return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + graph.Id)}, err
						} else {
							return models.ResponseType{Type: 500, Error: err}, err
						}
					}
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
				return models.ResponseType{Type: 500, Error: err}, err
			}

			err = cur.All(ctx, &addedGraphs)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}

			cur.Close(ctx)

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
					//return err, err.Error
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
	}, txnOptions)

	result := res.(models.ResponseType)
	return models.ResponseType{Type: result.Type, Error: result.Error}
}

func (s *MongoDB) DeleteFloor(id string) models.ResponseType {
	stairsCol := s.Database.Collection("stairs")
	graphsCol := s.Database.Collection("graph_points")
	floorsCol := s.Database.Collection("floors")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return models.ResponseType{Type: 500, Error: err}
	}

	defer session.EndSession(context.TODO())

	res, _ := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, err
		}
		floorFilter := bson.M{
			"_id": objId,
		}
		var floor models.Floor
		err = floorsCol.FindOne(ctx, floorFilter).Decode(&floor)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified id: " + id)}, err
			} else {
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}

		if len(floor.Graph) > 0 {
			var graphs []models.GraphPoint
			filter := bson.M{
				"_id": bson.M{
					"$in": floor.Graph,
				},
			}

			cur, err := graphsCol.Find(ctx, filter)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}

			err = cur.All(ctx, &graphs)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}

			cur.Close(ctx)

			for _, graph := range graphs {
				if graph.StairId != "" {
					stairFilter := bson.M{
						"stairPoint": graph.StairId,
					}
					var stair models.Stair
					err = stairsCol.FindOne(ctx, stairFilter).Decode(&stair)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}, err
						} else {
							return models.ResponseType{Type: 500, Error: err}, err
						}
					}
					log.Println(stair.Id)

					linkIndex := utils.GetIndex(stair.Links, graph.Id)
					newLinks := stair.Links
					if linkIndex != -1 {
						newLinks = append(stair.Links[:linkIndex], stair.Links[linkIndex+1:]...)
					}

					_, err = stairsCol.UpdateOne(ctx, stairFilter, bson.M{"$set": bson.M{"links": newLinks}})
					if err != nil {
						return models.ResponseType{Type: 500, Error: err}, err
					}
				}
			}

			_, err = graphsCol.DeleteMany(ctx, filter)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}

		_, err = floorsCol.DeleteOne(ctx, floorFilter)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, err
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	}, txnOptions)

	result := res.(models.ResponseType)
	return models.ResponseType{Type: result.Type, Error: result.Error}
}
