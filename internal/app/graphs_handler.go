package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetGraphHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	graphData, err := s.Store.GetGraph(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetGraph")
	}

	return c.JSON(graphData)
}

func (s *API) GetAllGraphsHandler(c *fiber.Ctx) error {
	graphData, err := s.Store.GetAllGraphs()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetAllGraphs")
	}

	return c.JSON(graphData)
}

func (s *API) DeleteGraphHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	err := s.Store.DeleteGraph(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in DeleteGraph")
	}

	return err
}
