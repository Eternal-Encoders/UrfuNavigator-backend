package app

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetInstituteHandler(c *fiber.Ctx) error {
	url, urlExist := c.Queries()["url"]

	if !urlExist {
		log.Println("Request Institute without url")
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain url query parameters")
	}

	instituteData, res := s.Store.GetInstitute(url)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	iconData, iconResp := s.Store.GetInstituteIconsByName([]string{instituteData.Icon})
	if iconResp.Error != nil {
		log.Println(iconResp)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetInstituteIcons")
	}
	if len(iconData) != 1 {
		log.Println("There is too many or no media with id")
		log.Println(iconData)
		return c.Status(fiber.StatusNotFound).SendString("Cannot find media by id")
	}

	response := models.InstituteGet{
		Id:              instituteData.Id.Hex(),
		Name:            instituteData.Name,
		DisplayableName: instituteData.DisplayableName,
		MinFloor:        instituteData.MinFloor,
		MaxFloor:        instituteData.MaxFloor,
		Url:             instituteData.Url,
		Latitude:        instituteData.Latitude,
		Longitude:       instituteData.Longitude,
		Icon:            utils.IconToIconResponse(iconData[0]),
	}

	return c.JSON(response)
}

func (s *API) GetAllInstitutesHandler(c *fiber.Ctx) error {
	institutesData, res := s.Store.GetAllInstitutes()
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	iconIds := []string{}
	for _, institute := range institutesData {
		iconIds = append(iconIds, institute.Icon)
	}

	iconsData, iconResp := s.Store.GetInstituteIcons(iconIds)
	if iconResp.Error != nil {
		log.Println(iconResp)
		return c.Status(res.Type).SendString(iconResp.Error.Error())
	}

	if len(iconsData) != len(institutesData) {
		log.Printf("IconsData length = %d and InstitutesData length = %d", len(iconsData), len(institutesData))
		return c.Status(fiber.StatusNotFound).SendString("For some of the institutes icons not founded")
	}

	response := []models.InstituteGet{}
	for i, institue := range institutesData {
		response = append(response, models.InstituteGet{
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
	data := new(models.InstitutePost)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something wrong with request body")
	}

	res := s.Store.PostInstitute(*data)
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully created")
}

func (s *API) DeleteInstituteHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	res := s.Store.DeleteInstitute(id)
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully deleted")
}

func (s *API) PutInstituteHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	data := new(models.InstitutePost)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something wrong with request body")
	}

	res := s.Store.UpdateInstitute(*data, id)
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully updated")
}
