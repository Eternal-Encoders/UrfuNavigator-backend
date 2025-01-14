package services

import (
	"UrfuNavigator-backend/internal/models"
	"context"
	"errors"
	"log"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

func (s *Services) GetIconImage(context context.Context, iconName string) ([]byte, models.ResponseType) {
	obj, err := s.ObjStore.GetFile(iconName)
	if err != nil {
		log.Println(err)
		return obj, models.ResponseType{Type: 400, Error: errors.New("Cannot get file from Object Storage")}
	}

	return obj, models.ResponseType{Type: 200, Error: nil}
}

func (s *Services) GetIcon(context context.Context, name string) (models.InstituteIconPost, models.ResponseType) {
	icon, res := s.Store.GetInstituteIconsById(context, []string{name})
	if res.Error != nil {
		log.Println(res.Error)
		return models.InstituteIconPost{}, res
	}
	if len(icon) != 1 {
		log.Println("There is too many or no media with id")
		return models.InstituteIconPost{}, models.ResponseType{Type: 404, Error: errors.New("Cannot find media by id")}
	}

	response := models.InstituteIconPost{
		Id:  icon[0].Id.Hex(),
		Url: icon[0].Url,
		Alt: icon[0].Alt,
	}

	return response, models.ResponseType{Type: 200, Error: nil}
}

func (s *Services) GetAllIcons(context context.Context) ([]models.InstituteIconPost, models.ResponseType) {
	icons, res := s.Store.GetAllInstituteIcons(context)
	if res.Error != nil {
		log.Println(res.Error)
		return []models.InstituteIconPost{}, res
	}

	response := []models.InstituteIconPost{}
	for _, icon := range icons {
		response = append(response, models.InstituteIconPost{
			Id:  icon.Id.Hex(),
			Url: icon.Url,
			Alt: icon.Alt,
		})
	}

	return response, models.ResponseType{Type: 200, Error: nil}
}

func (s *Services) PostIcon(context context.Context, file *multipart.FileHeader) models.ResponseType {
	url, name, err := s.ObjStore.PostFile(*file)
	if err != nil {
		log.Println(err)
		return models.ResponseType{Type: 500, Error: errors.New("Error occured while loading file in bucket")}
	}
	if name == "" {
		return models.ResponseType{Type: fiber.StatusConflict, Error: errors.New("There is a file in bucket with this name")}
	}

	res := s.Store.InsertInstituteIcon(context, models.InstituteIconGet{
		Url: url,
		Alt: name,
	})

	return res
}

func (s *Services) DeleteIcon(context context.Context, id string) models.ResponseType {
	name, res := s.Store.DeleteInstituteIcon(context, id)
	if res.Error != nil {
		return res
	}

	err := s.ObjStore.DeleteFile(name)
	if err != nil {
		log.Println(res)
		return models.ResponseType{Type: 500, Error: errors.New("Error occured while deleting file from bucket")}
	}

	return models.ResponseType{Type: 200, Error: nil}
}
