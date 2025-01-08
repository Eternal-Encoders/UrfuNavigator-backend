package app

import (
	"UrfuNavigator-backend/internal/models"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (s *API) RegisterHandler(c *fiber.Ctx) error {
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading body from request")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return c.Status(500).SendString("Something went wrong while generating hash from password")
	}

	hashedData := models.UserDB{
		Email: data.Email,
		Hash:  string(hashedPassword),
	}

	res := s.Store.Register(hashedData)
	if res.Error != nil {
		log.Println(err)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully registered")
}

func (s *API) LoginHandler(c *fiber.Ctx) error {
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading body from request")
	}

	user, res := s.Store.GetUser(data.Email)
	if res.Error != nil {
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(data.Password)); err != nil {
		return c.Status(406).SendString("wrong password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).JSON(t)
}

// func (s *API) LogoutHandler(c *fiber.Ctx) error {

// 	return c.Status(res.Type).SendString("successfully logged out")
// }
