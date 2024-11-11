package database

import "UrfuNavigator-backend/internal/models"

type Store interface {
	GetInstitute(url string) models.InstituteReadDBResponse
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
	PostInstituteIcon(models.InstituteIconRequest) error
	PostInstitute(models.InstituteRequest) error
	PostFloor(models.FloorFromFile) error
	PostGraphs(graphs []*models.GraphPoint) error
	PostStairs(graphs []*models.GraphPoint) error
	DeleteInstituteIcon(id string) (string, error)
	DeleteInstitute(id string) error
	DeleteFloor(id string) error
	DeleteGraph(id string) error
	DeleteStair(id string) error
	UpdateInstitute(body models.InstituteRequest, id string) error
	// UpdateInstituteIcon(body models.InstituteIconRequest, id string) error
}
