package app

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetInstituteHandler(c *fiber.Ctx) error {
	url := c.Query("url")

	if url == "" {
		log.Println("Request Institute without url")
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain url query parameters")
	}

	dbRes := s.Store.GetInstitute(url)
	if dbRes.Error != nil {
		log.Println(dbRes.Error)
		return c.Status(dbRes.ErrorType).SendString(dbRes.Error.Error())
	}

	iconData, iconErr := s.Store.GetInstituteIconsByName([]string{dbRes.Response.Icon})
	if iconErr != nil {
		log.Println(iconErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetInstituteIcons")
	}
	if len(iconData) != 1 {
		log.Println("There is too many or no media with id")
		log.Println(iconData)
		return c.Status(fiber.StatusNotFound).SendString("Cannot find media by id")
	}

	response := models.InstituteResponse{
		Id:              dbRes.Response.Id.Hex(),
		Name:            dbRes.Response.Name,
		DisplayableName: dbRes.Response.DisplayableName,
		MinFloor:        dbRes.Response.MinFloor,
		MaxFloor:        dbRes.Response.MaxFloor,
		Url:             dbRes.Response.Url,
		Latitude:        dbRes.Response.Latitude,
		Longitude:       dbRes.Response.Longitude,
		Icon:            utils.IconToIconResponse(iconData[0]),
	}

	return c.JSON(response)
}

func (s *API) GetAllInstitutesHandler(c *fiber.Ctx) error {
	institutesData, err := s.Store.GetAllInstitutes()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetInstitute")
	}

	iconIds := []string{}
	for _, institute := range institutesData {
		iconIds = append(iconIds, institute.Icon)
	}

	iconsData, err := s.Store.GetInstituteIcons(iconIds)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetInstituteIcons")
	}

	if len(iconsData) != len(institutesData) {
		log.Printf("IconsData length = %d and InstitutesData length = %d", len(iconsData), len(institutesData))
		return c.Status(fiber.StatusNotFound).SendString("For some of the institutes icons not founded")
	}

	response := []models.InstituteResponse{}
	for i, institue := range institutesData {
		response = append(response, models.InstituteResponse{
			Id:              institue.Id.Hex(),
			Name:            institue.Name,
			DisplayableName: institue.DisplayableName,
			MinFloor:        institue.MinFloor,
			MaxFloor:        institue.MaxFloor,
			Url:             institue.Url,
			Latitude:        institue.Latitude,
			Longitude:       institue.Longitude,
			Icon:            utils.IconToIconResponse(iconsData[i]),
		})
	}

	return c.JSON(response)
}

func (s *API) PostInstituteHandler(c *fiber.Ctx) error {
	data := new(models.InstituteRequest)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something wrong with request body")
	}

	err := s.Store.PostInstitute(*data)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something wrong in PostInstitute")
	}

	return err
}

func (s *API) DeleteInstituteHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	err := s.Store.DeleteInstitute(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in DeleteInstitute")
	}

	return err
}

func (s *API) PutInstituteHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	data := new(models.InstituteRequest)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something wrong with request body")
	}

	err := s.Store.UpdateInstitute(*data, id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in UpdateInstitute")
	}

	return err
}
