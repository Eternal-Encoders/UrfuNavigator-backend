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
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong while reading file into []bytes")
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

	floor := models.FloorRequest{
		Institute: floorFromFile.Institute,
		Floor:     floorFromFile.Floor,
		Width:     floorFromFile.Width,
		Height:    floorFromFile.Height,
		Service:   floorFromFile.Service,
		Audiences: audArr,
		Graph:     graphKeysArr,
	}

	res := s.Store.PostFloor(floor, graphArr)
	if res.Error != nil {
		log.Println(err)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	// if err == nil {
	// 	err := s.Store.PostGraphs(graphArr)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in PostGraphs")
	// 	} else {
	// 		err = s.Store.PostStairs(graphArr)
	// 		if err != nil {
	// 			log.Println(err)
	// 			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in PostStairs")
	// 		}
	// 	}
	// } else {

	// }

	return c.Status(res.Type).SendString("successfully created")
}

func (s *API) GetFloorHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	log.Println(id)

	floorData, res := s.Store.GetFloor(id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
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
	floorData, res := s.Store.GetAllFloors()
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
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

func (s *API) PutFloorHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	data := new(models.FloorPut)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something wrong with request body")
	}

	res := s.Store.UpdateFloor(*data, id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully updated")
}

func (s *API) DeleteFloorHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	res := s.Store.DeleteFloor(id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully deleted")
}
