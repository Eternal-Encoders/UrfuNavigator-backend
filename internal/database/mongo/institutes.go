package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *MongoDB) GetInstitute(url string) (models.Institute, error) {
	collection := s.Database.Collection("insitutes")
	log.Println(url)
	filter := bson.M{
		"url": url,
	}

	var result models.Institute
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (s *MongoDB) GetAllInstitutes() ([]models.Institute, error) {
	collection := s.Database.Collection("institutes")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var result []models.Institute
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, decodeErr
	}
	// log.Println(result[0].Id)

	return result, nil
}

func (s *MongoDB) PostInstitute(institute models.InstitutePost) error {
	collection := s.Database.Collection("institutes")
	// iconCol := s.Database.Collection("media")

	filter := bson.M{"name": institute.Name}
	// log.Println(filter)
	err := collection.FindOne(context.TODO(), filter).Err()
	if err == nil {
		return errors.New("institute already exists")
	}

	_, err = collection.InsertOne(context.TODO(), institute)
	return err
}

func (s *MongoDB) UpdateInstitute(body models.InstitutePost, id string) error {
	instituteCol := s.Database.Collection("insitutes")
	floorsCol := s.Database.Collection("floors")
	graphsCol := s.Database.Collection("graph_points")
	stairsCol := s.Database.Collection("stairs")
	mediaCol := s.Database.Collection("media")

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return err
	}

	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		filter := bson.M{
			"_id": objId,
		}

		if body.Name == "" {
			return nil, errors.New("institute name can not be empty")
		}

		var oldInstitute models.Institute
		err = instituteCol.FindOne(ctx, filter).Decode(&oldInstitute)
		if err != nil {
			return nil, err
		}

		if oldInstitute.Name != body.Name || oldInstitute.MaxFloor != body.MaxFloor || oldInstitute.MinFloor != body.MinFloor {
			var floors []models.Floor

			cur, err := floorsCol.Find(ctx, bson.M{"institute": oldInstitute.Name})
			if err != nil {
				return nil, err
			}

			err = cur.All(ctx, &floors)
			if err != nil {
				return nil, err
			}
			cur.Close(ctx)

			for _, floor := range floors {
				if body.MinFloor > floor.Floor || floor.Floor > body.MaxFloor {
					return nil, errors.New("there is a floor out of new floor bounds")
				}

				_, err = graphsCol.UpdateMany(ctx, bson.M{"institute": oldInstitute.Name, "floor": floor.Floor},
					bson.M{"$set": bson.M{"institute": body.Name}})
				if err != nil {
					return nil, err
				}
			}

			_, err = stairsCol.UpdateMany(ctx, bson.M{"institute": oldInstitute.Name}, bson.M{"$set": bson.M{"institute": body.Name}})
			if err != nil {
				return nil, err
			}

			_, err = floorsCol.UpdateMany(ctx, bson.M{"institute": oldInstitute.Name}, bson.M{"$set": bson.M{"institute": body.Name}})
			if err != nil {
				return nil, err
			}
		}

		if oldInstitute.Icon != body.Icon {
			err = mediaCol.FindOne(ctx, bson.M{"url": oldInstitute.Icon}).Err()
			if err != nil {
				return nil, err
			}
		}

		_, err = instituteCol.UpdateOne(ctx, filter, bson.M{"$set": body})
		fmt.Println(err != nil)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	return err
}

func (s *MongoDB) DeleteInstitute(id string) error {
	collection := s.Database.Collection("institutes")
	floorsCol := s.Database.Collection("floors")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": objId,
	}

	var institute models.Institute
	if err = collection.FindOne(context.TODO(), filter).Decode(&institute); err != nil {
		return err
	}

	if err = floorsCol.FindOne(context.TODO(), bson.M{"institute": institute.Name}).Err(); err == nil {
		return errors.New("can not delete institute with floors")
	}

	_, err = collection.DeleteOne(context.TODO(), filter)
	return err
}
