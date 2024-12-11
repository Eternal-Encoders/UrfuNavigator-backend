package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
)

type Store interface {
	GetInstitute(url string) (models.Institute, models.ResponseType)
	GetAllInstitutes() ([]models.Institute, models.ResponseType)
	GetInstituteIcons(ids []string) ([]models.InstituteIcon, models.ResponseType)
	GetInstituteIconsByName(ids []string) ([]models.InstituteIcon, models.ResponseType)
	GetAllInstituteIcons() ([]models.InstituteIcon, models.ResponseType)
	GetFloor(id string) (models.Floor, models.ResponseType)
	GetAllFloors() ([]models.Floor, models.ResponseType)
	GetGraph(id string) (models.GraphPoint, models.ResponseType)
	GetAllGraphs() ([]models.GraphPoint, models.ResponseType)
	GetStair(id string) (models.Stair, models.ResponseType)
	GetAllStairs() ([]models.Stair, models.ResponseType)
	PostInstituteIcon(models.InstituteIconGet) models.ResponseType
	PostInstitute(models.InstitutePost) models.ResponseType
	PostFloor(floor models.FloorRequest, graphs []*models.GraphPoint) models.ResponseType
	PostGraphs(context context.Context, graphs []*models.GraphPoint) models.ResponseType
	PostStairs(context context.Context, graphs []*models.GraphPoint) models.ResponseType
	DeleteInstituteIcon(id string) (string, models.ResponseType)
	DeleteInstitute(id string) models.ResponseType
	DeleteFloor(id string) models.ResponseType
	DeleteGraph(id string) models.ResponseType
	DeleteStair(id string) models.ResponseType
	UpdateInstitute(body models.InstitutePost, id string) models.ResponseType
	UpdateFloor(body models.FloorPut, id string) models.ResponseType
	UpdateGraph(context context.Context, body models.GraphPointPut, id string) models.ResponseType
	UpdateStair(context context.Context, body models.Stair, id string) models.ResponseType

	// UpdateInstituteIcon(body models.InstituteIconRequest, id string) error
}
