package mongo

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *MongoDB) Transaction(context context.Context, fn func(ctx mongo.SessionContext) (interface{}, error)) models.ResponseType {
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.Client.StartSession()
	if err != nil {
		log.Println("Something went wrong while starting new session")
		return models.ResponseType{Type: 500, Error: err}
	}

	defer session.EndSession(context)

	res, _ := session.WithTransaction(context, fn, txnOptions)
	result := res.(models.ResponseType)
	return models.ResponseType{Type: result.Type, Error: result.Error}
}
