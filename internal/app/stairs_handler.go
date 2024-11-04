package app

import (
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

func (s *API) DeleteStairHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	err := s.Store.DeleteStair(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetStair")
	}

	return err
}
