package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetUserHandler(c *fiber.Ctx) error {
	email := c.Query("email")

	userData, res := s.Store.GetUser(email)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(userData)
}

func (s *API) GetAllUsersHandler(c *fiber.Ctx) error {
	userData, res := s.Store.GetAllUsers()
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(userData)
}
