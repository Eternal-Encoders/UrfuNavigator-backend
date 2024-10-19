package database

import "UrfuNavigator-backend/internal/models"

type Store interface {
	GetInstitute(url string) (models.Institute, error)
	GetAllInstitutes() ([]models.Institute, error)
	GetInstituteIcons(ids []string) ([]models.InstituteIcon, error)
	GetInstituteIconsByName(ids []string) ([]models.InstituteIcon, error)
	GetAllInstituteIcons() ([]models.InstituteIcon, error)
	PostInstituteIcon(models.InstituteIconRequest) error
	PostInstitute(models.InstituteRequest) error
	DeleteInstituteIcon(id string) (string, error)
	DeleteInstitute(id string) error
	UpdateInstitute(body models.InstituteRequest, id string) error
	// UpdateInstituteIcon(body models.InstituteIconRequest, id string) error
}
