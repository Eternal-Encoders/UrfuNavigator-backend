package app

import (
	"context"
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

	obj, res := s.Services.GetIconImage(context.TODO(), iconName)
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	c.Attachment(iconName)
	return c.Status(200).Send(obj)
}

func (s *API) GetIconHandler(c *fiber.Ctx) error {
	iconName := c.Query("id")
	icon, res := s.Services.GetIcon(context.TODO(), iconName)
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.JSON(icon)
}

func (s *API) GetAllIconsHandler(c *fiber.Ctx) error {
	icons, res := s.Services.GetAllIcons(context.TODO())
	if res.Error != nil {
		log.Println(res.Error)
		return c.Status(res.Type).SendString(res.Error.Error())
	}

	return c.Status(res.Type).JSON(icons)
}

func (s *API) PostIconHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("icon")

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Error occured while loading file from request")
	}

	res := s.Services.PostIcon(context.TODO(), file)

	return c.Status(res.Type).SendString("successfully added")
}

func (s *API) DeleteIconHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	log.Println(id)
	res := s.Services.DeleteIcon(context.TODO(), id)

	return c.Status(res.Type).SendString("successfully deleted")
}
