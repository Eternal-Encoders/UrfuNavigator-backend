package app

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) PostFloorFromFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("floor")

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading file from request")
	}

	res := s.Services.PostFloor(context.TODO(), file)
	if res.Error != nil {
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully created")
}

func (s *API) GetFloorHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	log.Println(id)

	floorData, res := s.Services.GetFloor(context.TODO(), id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(200).JSON(floorData)
}

func (s *API) GetAllFloorsHandler(c *fiber.Ctx) error {
	floorData, res := s.Services.GetAllFloors(context.TODO())
	if res.Error != nil {
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(200).JSON(floorData)
}

func (s *API) PutFloorHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	data := new(models.FloorPut)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something wrong with request body")
	}

	res := s.Services.UpdateFloor(context.TODO(), *data, id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully updated")
}

func (s *API) DeleteFloorHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	res := s.Services.DeleteFloor(context.TODO(), id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully deleted")
}
