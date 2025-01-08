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

func (s *MongoDB) GetInstitute(url string) (models.Institute, models.ResponseType) {
	collection := s.Database.Collection("insitutes")
	filter := bson.M{
		"url": url,
	}

	var result models.Institute
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified url: " + url)}
		} else {
			return result, models.ResponseType{Type: 500, Error: err}
		}
	}

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) GetAllInstitutes() ([]models.Institute, models.ResponseType) {
	collection := s.Database.Collection("insitutes")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}

	defer cursor.Close(context.TODO())

	var result []models.Institute
	decodeErr := cursor.All(context.TODO(), &result)
	if decodeErr != nil {
		return nil, models.ResponseType{Type: 500, Error: err}
	}
	// if len(result) == 0 {
	// 	return nil, models.ResponseType{Type: 404, Error: errors.New("there are no institutes")}
	// }

	return result, models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) PostInstitute(institute models.InstitutePost) models.ResponseType {
	collection := s.Database.Collection("insitutes")
	// iconCol := s.Database.Collection("media")

	filter := bson.M{"name": institute.Name}

	err := collection.FindOne(context.TODO(), filter).Err()
	if err == nil {
		return models.ResponseType{Type: 406, Error: errors.New("institute already exists")}
	}

	_, err = collection.InsertOne(context.TODO(), institute)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) UpdateInstitute(body models.InstitutePost, id string) models.ResponseType {
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
		return models.ResponseType{Type: 500, Error: err}
	}

	defer session.EndSession(context.TODO())

	res, _ := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, err
		}
		filter := bson.M{
			"_id": objId,
		}

		if body.Name == "" {
			err = errors.New("institute name can not be empty")
			return models.ResponseType{Type: 406, Error: err}, err
		}

		var oldInstitute models.Institute
		err = instituteCol.FindOne(ctx, filter).Decode(&oldInstitute)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified id: " + id)}, err
			} else {
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}

		if oldInstitute.Name != body.Name || oldInstitute.MaxFloor != body.MaxFloor || oldInstitute.MinFloor != body.MinFloor {
			var floors []models.Floor

			cur, err := floorsCol.Find(ctx, bson.M{"institute": oldInstitute.Name})
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}

			err = cur.All(ctx, &floors)
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}
			cur.Close(ctx)

			// if len(floors) == 0 {
			// 	return models.ResponseType{Type: 404, Error: errors.New("there are no floors with specified institute: " + oldInstitute.Name)}, err
			// }

			for _, floor := range floors {
				if body.MinFloor > floor.Floor || floor.Floor > body.MaxFloor {
					err = errors.New("there is a floor out of new floor bounds")
					return models.ResponseType{Type: 406, Error: err}, err
				}

				_, err = graphsCol.UpdateMany(ctx, bson.M{"institute": oldInstitute.Name, "floor": floor.Floor},
					bson.M{"$set": bson.M{"institute": body.Name}})
				if err != nil {
					return models.ResponseType{Type: 500, Error: err}, err
				}
			}

			_, err = stairsCol.UpdateMany(ctx, bson.M{"institute": oldInstitute.Name}, bson.M{"$set": bson.M{"institute": body.Name}})
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}

			_, err = floorsCol.UpdateMany(ctx, bson.M{"institute": oldInstitute.Name}, bson.M{"$set": bson.M{"institute": body.Name}})
			if err != nil {
				return models.ResponseType{Type: 500, Error: err}, err
			}
		}

		if oldInstitute.Icon != body.Icon {
			err = mediaCol.FindOne(ctx, bson.M{"url": oldInstitute.Icon}).Err()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return models.ResponseType{Type: 404, Error: errors.New("there is no icon with specified name: " + oldInstitute.Icon)}, err
				} else {
					return models.ResponseType{Type: 500, Error: err}, err
				}
			}
		}

		_, err = instituteCol.UpdateOne(ctx, filter, bson.M{"$set": body})
		fmt.Println(err != nil)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, err
		}

		return models.ResponseType{Type: 200, Error: nil}, nil
	}, txnOptions)

	result := res.(models.ResponseType)
	return models.ResponseType{Type: result.Type, Error: result.Error}
}

func (s *MongoDB) DeleteInstitute(id string) models.ResponseType {
	collection := s.Database.Collection("insitutes")
	floorsCol := s.Database.Collection("floors")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}
	filter := bson.M{
		"_id": objId,
	}

	var institute models.Institute
	if err = collection.FindOne(context.TODO(), filter).Decode(&institute); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified id: " + id)}
		} else {
			return models.ResponseType{Type: 500, Error: err}
		}
	}

	if err = floorsCol.FindOne(context.TODO(), bson.M{"institute": institute.Name}).Err(); err == nil {
		if err == mongo.ErrNoDocuments {
			return models.ResponseType{Type: 404, Error: errors.New("there is no floor with specified institute: " + institute.Name)}
		} else {
			return models.ResponseType{Type: 500, Error: err}
		}
	}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}
	return models.ResponseType{Type: 200, Error: nil}
}
