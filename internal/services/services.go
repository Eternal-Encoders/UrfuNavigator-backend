package services

import (
	"UrfuNavigator-backend/internal/database"
	"UrfuNavigator-backend/internal/objstore"
)

type Services struct {
	Store    database.Store
	ObjStore objstore.ObjectStore
}

// type Services struct {
// 	UserServices      UserServices
// 	InstituteServices InstituteServices
// 	MediaServices     MediaServices
// 	FloorServices     FloorServices
// 	GraphServices     GraphServices
// 	StairServices     StairServices
// }

// func NewServices(s database.Store, o objstore.ObjectStore) Services {
// 	return Services{
// 		UserServices:      UserServices{Store: s},
// 		InstituteServices: InstituteServices{Store: s},
// 		MediaServices:     MediaServices{Store: s, ObjStore: o},
// 		FloorServices:     FloorServices{Store: s},
// 		GraphServices:     GraphServices{Store: s},
// 		StairServices:     StairServices{Store: s},
// 	}
// }

// GetFloor(id string) (models.Floor, models.ResponseType)
// GetAllFloors() ([]models.Floor, models.ResponseType)
// PostFloor(floor models.FloorRequest, graphs []*models.GraphPoint) models.ResponseType
// DeleteFloor(id string) models.ResponseType
// UpdateFloor(body models.FloorPut, id string) models.ResponseType

// GetGraph(id string) (models.GraphPoint, models.ResponseType)
// GetAllGraphs() ([]models.GraphPoint, models.ResponseType)
// PostGraphs(context context.Context, graphs []*models.GraphPoint) models.ResponseType
// DeleteGraph(id string) models.ResponseType
// UpdateGraph(context context.Context, body models.GraphPointPut, id string) models.ResponseType

// GetStair(id string) (models.Stair, models.ResponseType)
// GetAllStairs() ([]models.Stair, models.ResponseType)
// PostStairs(context context.Context, graphs []*models.GraphPoint) models.ResponseType
// DeleteStair(id string) models.ResponseType
// UpdateStair(context context.Context, body models.Stair, id string) models.ResponseType
