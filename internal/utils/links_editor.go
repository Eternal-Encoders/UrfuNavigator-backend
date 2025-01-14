package utils

import (
	"UrfuNavigator-backend/internal/database"
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"
	"log"
	"slices"
	"time"
)

func GetAddedDeleted(oldSlice []string, newSlice []string) ([]string, []string) {
	var added []string
	var deleted []string

	for _, v := range oldSlice {
		if !slices.Contains(newSlice, v) {
			deleted = append(deleted, v)
		}
	}

	for _, v := range newSlice {
		if !slices.Contains(oldSlice, v) {
			added = append(added, v)
		}
	}

	return added, deleted
}

func UpdateGraphsStairs(ctx context.Context, graphsId []string, stair models.Stair, Store database.Store) models.ResponseType {
	var graphs []models.GraphPoint
	start := time.Now()

	if len(graphsId) == 0 {
		return models.ResponseType{Type: 200, Error: nil}
	}

	// filter := bson.M{
	// 	"_id": bson.M{
	// 		"$in": graphsId,
	// 	},
	// }

	graphs, res := Store.GetManyGraphsByIds(ctx, graphsId)
	if res.Error != nil {
		return res
	}

	// cur, err := graphsCol.Find(ctx, filter)
	// if err != nil {
	// 	return models.ResponseType{Type: 500, Error: err}
	// }

	// err = cur.All(ctx, &graphs)
	// if err != nil {
	// 	return models.ResponseType{Type: 500, Error: err}
	// }

	// cur.Close(ctx)

	if len(graphs) < len(graphsId) {
		return models.ResponseType{Type: 404, Error: errors.New("some graph point is missing in database")}
	}

	for _, graph := range graphs {
		// for _, id := range graphsId {
		// var graph models.GraphPoint
		// err := graphsCol.FindOne(ctx, bson.M{"_id": id}).Decode(&graph)
		// if err != nil {
		// 	return err
		// //Здесь можно указать, какая конкретно точка отсутствует
		// }

		if graph.Institute != stair.Institute {
			return models.ResponseType{Type: 406, Error: errors.New("graph point and stair have different institute values")}
		}

		if slices.Contains(graph.Types, "stair") {
			var oldStair models.Stair

			filter := []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: graph.StairId}}
			oldStair, res = Store.GetStair(ctx, filter)
			if res.Error != nil {
				return res
			}

			// err = stairsCol.FindOne(ctx, bson.M{"_id": graph.StairId}).Decode(&oldStair)
			// if err != nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}
			// 	}
			// }

			graphIndex := GetIndex(oldStair.Links, graph.Id)
			newLinks := oldStair.Links
			if graphIndex != -1 {
				newLinks = append(oldStair.Links[:graphIndex], oldStair.Links[graphIndex+1:]...)
			}

			oldStair.Links = newLinks

			filter = []models.Query{models.Query{ParamName: "_id", Type: "string", StringValue: oldStair.Id}}
			res = Store.UpdateStair(ctx, filter, oldStair)
			if res.Error != nil {
				return res
			}

			// _, err = stairsCol.UpdateOne(context.TODO(), bson.M{"_id": oldStair.Id}, bson.M{"$set": bson.M{"links": newLinks}})
			// if err != nil {
			// 	if err == mongo.ErrNoDocuments {
			// 		return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + oldStair.Id)}
			// 	} else {
			// 		return models.ResponseType{Type: 500, Error: err}
			// 	}
			// }
		}

		// *graph.StairId = stair.Id
		if !slices.Contains(graph.Types, "stair") {
			graph.Types = append(graph.Types, "stair")
		}

		if !slices.Contains(graph.Types, "stair") {
			graph.Types = append(graph.Types, "stair")
		}

		filter := []models.Query{{ParamName: "_id", Type: "string", StringValue: graph.Id}}
		res = Store.UpdateGraph(ctx, filter, GraphPointToGraphPointPut(graph))
		if res.Error != nil {
			return res
		}

		// _, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graph.Id}, bson.M{"$set": graph})
		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + graph.Id)}
		// 	} else {
		// 		return models.ResponseType{Type: 500, Error: err}
		// 	}
		// }
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.Println(elapsed.Nanoseconds())

	return models.ResponseType{Type: 200, Error: nil}
}
