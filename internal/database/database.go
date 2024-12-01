package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"
)

type Store interface {
	GetInstitute(url string) (models.Institute, error)
	GetAllInstitutes() ([]models.Institute, error)
	GetInstituteIcons(ids []string) ([]models.InstituteIcon, error)
	GetInstituteIconsByName(ids []string) ([]models.InstituteIcon, error)
	GetAllInstituteIcons() ([]models.InstituteIcon, error)
	GetFloor(id string) (models.Floor, error)
	GetAllFloors() ([]models.Floor, error)
	GetGraph(id string) (models.GraphPoint, error)
	GetAllGraphs() ([]models.GraphPoint, error)
	GetStair(id string) (models.Stair, error)
	GetAllStairs() ([]models.Stair, error)
	PostInstituteIcon(models.InstituteIconGet) error
	PostInstitute(models.InstitutePost) error
	PostFloor(floor models.FloorRequest, graphs []*models.GraphPoint) error
	PostGraphs(context context.Context, graphs []*models.GraphPoint) error
	PostStairs(context context.Context, graphs []*models.GraphPoint) error
	DeleteInstituteIcon(id string) (string, error)
	DeleteInstitute(id string) error
	DeleteFloor(id string) error
	DeleteGraph(id string) error
	DeleteStair(id string) error
	UpdateInstitute(body models.InstitutePost, id string) error
	UpdateFloor(body models.FloorPut, id string) error
	UpdateGraph(context context.Context, body models.GraphPointPut, id string) error
	UpdateStair(context context.Context, body models.Stair, id string) error

	// UpdateInstituteIcon(body models.InstituteIconRequest, id string) error
}
