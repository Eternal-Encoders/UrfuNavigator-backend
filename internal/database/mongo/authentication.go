package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (s *MongoDB) Register(body models.UserDB) models.ResponseType {
	collection := s.Database.Collection("users")

	filter := bson.M{
		"url": body.Email,
	}

	err := collection.FindOne(context.TODO(), filter).Err()
	if err == nil {
		if err == models.ErrNoDocuments {
			return models.ResponseType{Type: 406, Error: errors.New("user with specified email already exists")}
		} else {
			return models.ResponseType{Type: 500, Error: err}
		}
	}

	_, err = collection.InsertOne(context.TODO(), body)
	if err != nil {
		return models.ResponseType{Type: 500, Error: err}
	}

	return models.ResponseType{Type: 200, Error: nil}
}

func (s *MongoDB) Login(body models.User) (models.ResponseType, *jwt.Token) {
	collection := s.Database.Collection("users")

	filter := bson.M{
		"email": body.Email,
	}

	var user models.UserDB

	if err := collection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		if err == models.ErrNoDocuments {
			return models.ResponseType{Type: 404, Error: errors.New("there is no user with specified email: " + body.Email)}, nil
		} else {
			return models.ResponseType{Type: 500, Error: err}, nil
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(body.Password)); err != nil {
		return models.ResponseType{Type: 406, Error: errors.New("wrong password")}, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    body.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(),
	})

	return models.ResponseType{Type: 200, Error: nil}, token
}

// func (s *MongoDB) Logout() models.ResponseType {
// 	collection := s.Database.Collection("users")

//  return models.ResponseType{Type: 200, Error: nil}
// }
