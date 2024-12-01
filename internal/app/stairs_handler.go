package app

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetStairHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	stairData, err := s.Store.GetStair(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetStair")
	}

	return c.JSON(stairData)
}

func (s *API) GetAllStairsHandler(c *fiber.Ctx) error {
	stairData, err := s.Store.GetAllStairs()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetAllStairs")
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

	err := s.Store.UpdateStair(context.TODO(), *data, id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in UpdateStair")
	}

	return err
}

func (s *API) DeleteStairHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	err := s.Store.DeleteStair(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetStair")
	}

	return err
}
