package database

import (
	"UrfuNavigator-backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	GetInstitute(context context.Context, filter []models.Query) (models.Institute, models.ResponseType)
	GetManyInstitutes(context context.Context, filter []models.Query) ([]models.Institute, models.ResponseType)
	GetAllInstitutes(context context.Context) ([]models.Institute, models.ResponseType)
	InsertInstitute(context context.Context, institute models.InstitutePost) models.ResponseType
	// InsertManyInstitutes(context context.Context, institutes []models.InstitutePost) models.ResponseTypes
	DeleteInstitute(context context.Context, filter []models.Query) models.ResponseType
	DeleteManyInstitutes(context context.Context, filter []models.Query) models.ResponseType
	UpdateInstitute(context context.Context, filter []models.Query, body models.InstitutePost) models.ResponseType
	UpdateManyInstitutes(context context.Context, filter []models.Query, body interface{}) models.ResponseType

	GetInstituteIconsById(context context.Context, ids []string) ([]models.InstituteIcon, models.ResponseType)
	GetInstituteIconsByName(context context.Context, names []string) ([]models.InstituteIcon, models.ResponseType)
	GetAllInstituteIcons(context context.Context) ([]models.InstituteIcon, models.ResponseType)
	InsertInstituteIcon(context context.Context, institute models.InstituteIconGet) models.ResponseType
	DeleteInstituteIcon(context context.Context, id string) (string, models.ResponseType)

	GetFloor(context context.Context, filter []models.Query) (models.Floor, models.ResponseType)
	GetManyFloors(context context.Context, filter []models.Query) ([]models.Floor, models.ResponseType)
	GetAllFloors(context context.Context) ([]models.Floor, models.ResponseType)
	InsertFloor(context context.Context, floor models.FloorRequest) models.ResponseType
	// InsertManyFloors(context context.Context, floors []models.FloorRequest) models.ResponseType
	DeleteFloor(context context.Context, filter []models.Query) models.ResponseType
	DeleteManyFloors(context context.Context, filter []models.Query) models.ResponseType
	UpdateFloor(context context.Context, filter []models.Query, body models.FloorPut) models.ResponseType
	UpdateManyFloors(context context.Context, filter []models.Query, body interface{}) models.ResponseType

	GetGraph(context context.Context, filter []models.Query) (models.GraphPoint, models.ResponseType)
	GetManyGraphs(context context.Context, filter []models.Query) ([]models.GraphPoint, models.ResponseType)
	GetAllGraphs(context context.Context) ([]models.GraphPoint, models.ResponseType)
	InsertManyGraphs(context context.Context, graphs []*models.GraphPoint) models.ResponseType
	GetManyGraphsByIds(context context.Context, ids []string) ([]models.GraphPoint, models.ResponseType)
	DeleteGraph(context context.Context, filter []models.Query) models.ResponseType
	DeleteManyGraphs(context context.Context, filter []models.Query) models.ResponseType
	DeleteManyGraphsByIds(context context.Context, ids []string) models.ResponseType
	UpdateGraph(context context.Context, filter []models.Query, body models.GraphPointPut) models.ResponseType
	UpdateManyGraphs(context context.Context, filter []models.Query, body interface{}) models.ResponseType

	GetStair(context context.Context, filter []models.Query) (models.Stair, models.ResponseType)
	GetManyStairs(context context.Context, filter []models.Query) ([]models.Stair, models.ResponseType)
	GetAllStairs(context context.Context) ([]models.Stair, models.ResponseType)
	InsertManyStairs(context context.Context, graphs []*models.Stair) models.ResponseType
	DeleteStair(context context.Context, filter []models.Query) models.ResponseType
	DeleteManyStairs(context context.Context, filter []models.Query) models.ResponseType
	UpdateStair(context context.Context, filter []models.Query, body models.Stair) models.ResponseType
	UpdateManyStairs(context context.Context, filter []models.Query, body interface{}) models.ResponseType

	GetUser(email string) (models.UserDB, models.ResponseType)
	GetAllUsers() ([]models.UserDB, models.ResponseType)
	InsertUser(user models.UserCreate) models.ResponseType

	////////
	Transaction(context context.Context, fn func(ctx mongo.SessionContext) (interface{}, error)) models.ResponseType
}
