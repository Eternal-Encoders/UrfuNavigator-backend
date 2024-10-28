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
	var floorFromFile models.FloorFromFile

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
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong while reading file into []bytes")
	}

	json.Unmarshal([]byte(buf.Bytes()), &floorFromFile)

	audArr := []*models.Auditorium{}
	for _, v := range floorFromFile.Audiences {
		audArr = append(audArr, v)
	}

	graphArr := []*models.GraphPoint{}
	graphKeysArr := []string{}
	for k, v := range floorFromFile.Graph {
		graphArr = append(graphArr, v)
		graphKeysArr = append(graphKeysArr, k)
	}

	floor := models.Floor{
		Institute: floorFromFile.Institute,
		Floor:     floorFromFile.Floor,
		Width:     floorFromFile.Width,
		Height:    floorFromFile.Height,
		Service:   floorFromFile.Service,
		Audiences: audArr,
		Graph:     graphKeysArr,
	}

	err = s.Store.PostFloor(floor)
	if err == nil {
		err := s.Store.PostGraphs(graphArr)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in PostGraphs")
		} else {
			err = s.Store.PostStairs(graphArr)
			if err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in PostStairs")
			}
		}
	} else {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in PostFloor")
	}

	// log.Println(floor.Audiences)
	return err
}
