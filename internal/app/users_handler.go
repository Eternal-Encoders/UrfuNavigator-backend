package app

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetUserHandler(c *fiber.Ctx) error {
	email := c.Query("email")

	userData, res := s.Services.GetUser(context.TODO(), email)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(userData)
}

func (s *API) GetAllUsersHandler(c *fiber.Ctx) error {
	userData, res := s.Services.GetAllUsers(context.TODO())
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(userData)
}
