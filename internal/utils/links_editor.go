package utils

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"
	"log"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func UpdateGraphsStairs(ctx context.Context, graphsId []string, stair models.Stair, graphsCol *mongo.Collection, stairsCol *mongo.Collection) models.ResponseType {
	var graphs []models.GraphPoint
	start := time.Now()

	if len(graphsId) == 0 {
		return models.ResponseType{Type: 200, Error: nil}
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": graphsId,
		},
	}

	cur, err := graphsCol.Find(ctx, filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	err = cur.All(ctx, &graphs)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	cur.Close(ctx)

	if len(graphs) < len(graphsId) {
		err = errors.New("some graph point is missing in database")
		return models.ResponseType{Type: 404, Error: err}
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
			err = errors.New("graph point and stair have different institute values")
			return models.ResponseType{Type: 406, Error: err}
		}

		if slices.Contains(graph.Types, "stair") {
			var oldStair models.Stair
			err = stairsCol.FindOne(ctx, bson.M{"_id": graph.StairId}).Decode(&oldStair)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + graph.StairId)}
				} else {
					return models.ResponseType{Type: 500, Error: err}
				}
			}

			graphIndex := GetIndex(oldStair.Links, graph.Id)
			newLinks := oldStair.Links
			if graphIndex != -1 {
				newLinks = append(oldStair.Links[:graphIndex], oldStair.Links[graphIndex+1:]...)
			}

			_, err = stairsCol.UpdateOne(context.TODO(), bson.M{"_id": oldStair.Id}, bson.M{"$set": bson.M{"links": newLinks}})
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return models.ResponseType{Type: 404, Error: errors.New("there is no stair with specified id: " + oldStair.Id)}
				} else {
					return models.ResponseType{Type: 500, Error: err}
				}
			}
		}

		// *graph.StairId = stair.Id
		if !slices.Contains(graph.Types, "stair") {
			graph.Types = append(graph.Types, "stair")
		}

		if !slices.Contains(graph.Types, "stair") {
			graph.Types = append(graph.Types, "stair")
		}

		_, err = graphsCol.UpdateOne(ctx, bson.M{"_id": graph.Id}, bson.M{"$set": graph})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no graph point with specified id: " + graph.Id)}
			} else {
				return models.ResponseType{Type: 500, Error: err}
			}
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.Println(elapsed.Nanoseconds())

	return models.ResponseType{Type: 200, Error: nil}
}
