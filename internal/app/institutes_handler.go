package app

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetInstituteHandler(c *fiber.Ctx) error {
	url, urlExist := c.Queries()["url"]

	if !urlExist {
		log.Println("Request Institute without url")
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain url query parameters")
	}

	instituteData, res := s.Services.GetInstitute(context.TODO(), url)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(instituteData)
}

func (s *API) GetAllInstitutesHandler(c *fiber.Ctx) error {
	institutesData, res := s.Services.GetAllInstitutes(context.TODO())
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).JSON(institutesData)
}

func (s *API) PostInstituteHandler(c *fiber.Ctx) error {
	data := new(models.InstitutePost)

	if err := c.BodyParser(data); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something wrong with request body")
	}

	res := s.Services.PostInstitute(context.TODO(), *data)
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully created")
}

func (s *API) DeleteInstituteHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	res := s.Services.DeleteInstitute(context.TODO(), id)
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

	res := s.Services.UpdateInstitute(context.TODO(), *data, id)
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).SendString("successfully updated")
}
