package app

import (
	"UrfuNavigator-backend/internal/models"
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) PostFloorFromFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("floor")
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading file from request")
	}

	f, err := file.Open()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong while reading file")
	}

	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong while reading file into []byte")
	}

	var floorFromFile models.FloorFromFile
	json.Unmarshal([]byte(buf.Bytes()), &floorFromFile)

	err = s.Store.PostFloor(floorFromFile)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in PostFloor")
	}
	return err
}

func (s *API) GetFloorHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	log.Println(id)

	floorData, err := s.Store.GetFloor(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetFloor")
	}

	response := models.FloorResponse{
		Id:        floorData.Id.Hex(),
		Institute: floorData.Institute,
		Floor:     floorData.Floor,
		Width:     floorData.Width,
		Height:    floorData.Height,
		Audiences: floorData.Audiences,
		Service:   floorData.Service,
		Graph:     floorData.Graph,
	}

	return c.JSON(response)
}

func (s *API) GetAllFloorsHandler(c *fiber.Ctx) error {
	floorData, err := s.Store.GetAllFloors()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetAllFloors")
	}

	response := []models.FloorResponse{}
	for _, floor := range floorData {
		response = append(response, models.FloorResponse{
			Id:        floor.Id.Hex(),
			Institute: floor.Institute,
			Floor:     floor.Floor,
			Width:     floor.Width,
			Height:    floor.Height,
			Audiences: floor.Audiences,
			Service:   floor.Service,
			Graph:     floor.Graph,
		})
	}

	return c.JSON(response)
}

// func (s *API) UpdateFloor(c *fiber.Ctx) error {

// }

func (s *API) DeleteFloorHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	err := s.Store.DeleteFloor(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in DeleteFloor")
	}

	return err
}
