package app

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) RegisterHandler(c *fiber.Ctx) error {
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading body from request")
	}

	res := s.Services.Register(context.TODO(), data)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully registered")
}

func (s *API) LoginHandler(c *fiber.Ctx) error {
	var data models.User

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading body from request")
	}

	token, res := s.Services.Login(context.TODO(), data)
	if res.Error != nil {
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).JSON(token)
}

// func (s *API) LogoutHandler(c *fiber.Ctx) error {

// 	return c.Status(res.Type).SendString("successfully logged out")
// }
