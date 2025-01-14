package app

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetGraphHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	graphData, res := s.Services.GetGraph(context.TODO(), id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(graphData)
}

func (s *API) GetAllGraphsHandler(c *fiber.Ctx) error {
	graphData, res := s.Services.GetAllGraphs(context.TODO())
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(graphData)
}

func (s *API) PutGraphHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	data := new(models.GraphPointPut)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something wrong with request body")
	}

	res := s.Services.UpdateGraph(context.TODO(), *data, id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully updated")
}

func (s *API) DeleteGraphHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	res := s.Services.DeleteGraph(context.TODO(), id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully deleted")
}
