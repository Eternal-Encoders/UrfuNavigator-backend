package services

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (s *Services) GetUser(context context.Context, email string) (models.UserDTO, models.ResponseType) {
	var userDTO models.UserDTO

	user, res := s.Store.GetUser(email)
	if res.Error != nil {
		return userDTO, res
	}

	userDTO.Email = user.Email
	userDTO.Id = user.Id.Hex()

	return userDTO, res
}

func (s *Services) GetAllUsers(context context.Context) ([]models.UserDTO, models.ResponseType) {
	var usersDTO []models.UserDTO

	users, res := s.Store.GetAllUsers()
	if res.Error != nil {
		return usersDTO, res
	}

	for _, v := range users {
		usersDTO = append(usersDTO, models.UserDTO{Email: v.Email, Id: v.Id.Hex()})
	}

	return usersDTO, res
}

func (s *Services) Register(context context.Context, body models.User) models.ResponseType {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		return models.ResponseType{Type: 500, Error: errors.New("Something went wrong while generating hash from password")}
	}

	hashedData := models.UserCreate{
		Email: body.Email,
		Hash:  string(hashedPassword),
	}

	res := s.Store.InsertUser(hashedData)
	return res
}

func (s *Services) Login(context context.Context, body models.User) (string, models.ResponseType) {
	user, res := s.Store.GetUser(body.Email)
	if res.Error != nil {
		return "", res
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(body.Password)); err != nil {
		return "", models.ResponseType{Type: 406, Error: errors.New("wrong password")}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", models.ResponseType{Type: 500, Error: err}
	}

	return t, models.ResponseType{Type: 200, Error: nil}
}

// func (s *MongoDB) Logout() models.ResponseType {
// 	collection := s.Database.Collection("users")

//  return models.ResponseType{Type: 200, Error: nil}
// }
