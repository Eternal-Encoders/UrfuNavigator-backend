package app

import (
	"UrfuNavigator-backend/internal/models"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (s *API) GetIconImageHandler(c *fiber.Ctx) error {
	iconName := c.Query("name")
	if !strings.HasSuffix(iconName, ".svg") {
		log.Println("Request Object with unsupported type")
		return c.Status(fiber.StatusBadRequest).SendString("This file type is not supported")
	}

	obj, err := s.ObjectStore.GetFile(iconName)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Cannot get file from Object Storage")
	}

	c.Attachment(iconName)
	return c.Send(obj)
}

func (s *API) GetIconHandler(c *fiber.Ctx) error {
	iconName := c.Query("id")
	icon, res := s.Store.GetInstituteIcons([]string{iconName})
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}
	if len(icon) != 1 {
		log.Println("There is too many or no media with id")
		return c.Status(fiber.StatusNotFound).SendString("Cannot find media by id")
	}

	response := models.InstituteIconPost{
		Id:  icon[0].Id.Hex(),
		Url: icon[0].Url,
		Alt: icon[0].Alt,
	}

	return c.JSON(response)
}

func (s *API) GetAllIconsHandler(c *fiber.Ctx) error {
	icons, res := s.Store.GetAllInstituteIcons()
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	response := []models.InstituteIconPost{}
	for _, icon := range icons {
		response = append(response, models.InstituteIconPost{
			Id:  icon.Id.Hex(),
			Url: icon.Url,
			Alt: icon.Alt,
		})
	}

	return c.JSON(response)
}

func (s *API) PostIconHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("icon")

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading file from request")
	}

	url, name, err := s.ObjectStore.PostFile(*file)
	fmt.Println(err != nil)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading file in bucket")
	}
	if name == "" {
		return c.Status(fiber.StatusConflict).SendString("There is a file in bucket with this name")
	}

	res := s.Store.PostInstituteIcon(models.InstituteIconGet{
		Url: url,
		Alt: name,
	})
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(200).SendString("successfully loaded")
}

func (s *API) DeleteIconHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	log.Println(id)
	name, res := s.Store.DeleteInstituteIcon(id)
	if res.Error != nil {
		log.Println(res)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	err := s.ObjectStore.DeleteFile(name)
	if err != nil {
		log.Println(res)
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while deleting file from bucket")
	}

	return c.Status(200).SendString("successfully deleted")
}
