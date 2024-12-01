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
	icon, err := s.Store.GetInstituteIcons([]string{iconName})
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in GetInstituteIcons")
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
	icons, err := s.Store.GetAllInstituteIcons()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in GetAllInstituteIcons")
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

	err = s.Store.PostInstituteIcon(models.InstituteIconGet{
		Url: url,
		Alt: name,
	})
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in PostInstituteIcon")
	}

	return err
}

func (s *API) DeleteIconHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	log.Println(id)
	name, err := s.Store.DeleteInstituteIcon(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in DeleteInstituteIcon")
	}

	err = s.ObjectStore.DeleteFile(name)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while deleting file from bucket")
	}
	return err
}
