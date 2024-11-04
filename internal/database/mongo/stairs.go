package database

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *MongoDB) PostStairs(graphs []*models.GraphPoint) error {
	collection := s.Database.Collection("stairs")

	if len(graphs) == 0 {
		return nil
	}

	for i := range graphs {
		if graphs[i].StairId != "" {
			filter := bson.D{{"_id", graphs[i].StairId}}
			var stairGraph models.Stair
			err := collection.FindOne(context.TODO(), filter).Decode(&stairGraph)
			if err != nil {
				_, err = collection.InsertOne(context.TODO(), models.Stair{
					Id:         graphs[i].StairId,
					StairPoint: graphs[i].StairId,
					Institute:  graphs[i].Institute,
					Links:      []string{graphs[i].Id},
				})

				if err != nil {
					return err
				}
			} else {
				l := append(stairGraph.Links, graphs[i].Id)
				_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"links": l}})
			}
		}
	}

	return nil
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

func (s *MongoDB) DeleteStair(id string) error {
	stairsCol := s.Database.Collection("stairs")
	graphsCol := s.Database.Collection("graph_points")

	stairFilter := bson.M{
		"_id": id,
	}
	var stair models.Stair
	err := stairsCol.FindOne(context.TODO(), stairFilter).Decode(&stair)
	if err != nil {
		return err
	}
	log.Println(stair.StairPoint)
	//////////////////////
	graphFilter := bson.M{
		"_id": bson.M{"$in": stair.Links},
	}
	cursor, err := graphsCol.Find(context.TODO(), graphFilter)
	if err != nil {
		return err
	}
	////////////////////
	defer cursor.Close(context.TODO())

	var graphs []models.GraphPoint
	decodeErr := cursor.All(context.TODO(), &graphs)
	if decodeErr != nil {
		return decodeErr
	}
	log.Println(graphs)
	///////////////////

	// _, err = graphsCol.UpdateMany(context.TODO(), bson.M{"_id": bson.M{"$in":ids}}, bson.M{"$set": bson.M{"stairId": ""}})
	// if err != nil {
	// 	return err
	// }

	for _, v := range graphs {
		log.Println(v)
		// linkIndex := utils.GetIndex(v.Links, stair.StairPoint)
		// log.Println(linkIndex)
		typeIndex := utils.GetIndex(v.Types, "stair")
		log.Println(typeIndex)
		newTypes := append(v.Types[:typeIndex], v.Types[typeIndex+1:]...)
		_, err = graphsCol.UpdateOne(context.TODO(), bson.M{"_id": v.Id}, bson.M{"$set": bson.M{"stairId": "", "types": newTypes}})
		if err != nil {
			return err
		}
	}

	_, err = stairsCol.DeleteOne(context.TODO(), stairFilter)
	if err != nil {
		return err
	}

	return err
}
