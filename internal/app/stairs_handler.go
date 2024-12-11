package app

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetStairHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	stairData, res := s.Store.GetStair(id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(stairData)
}

func (s *API) GetAllStairsHandler(c *fiber.Ctx) error {
	stairData, res := s.Store.GetAllStairs()
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(stairData)
}

func (s *API) PutStairHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	data := new(models.Stair)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something wrong with request body")
	}

	res := s.Store.UpdateStair(context.TODO(), *data, id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully updated")
}

func (s *API) DeleteStairHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	res := s.Store.DeleteStair(id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully deleted")
}
