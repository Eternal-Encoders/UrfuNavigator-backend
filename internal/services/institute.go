package services

import (
	"UrfuNavigator-backend/internal/models"
	"UrfuNavigator-backend/internal/utils"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Services) GetInstitute(context context.Context, url string) (models.InstituteGet, models.ResponseType) {
	var response models.InstituteGet
	institute, res := s.Store.GetInstitute(context,
		[]models.Query{{ParamName: "url", Type: "string", StringValue: url}})
	if res.Error != nil {
		return response, res
	}

	icon, res := s.Store.GetInstituteIconsByName(context, []string{institute.Icon})
	if res.Error != nil {
		return response, res
	}

	response = models.InstituteGet{
		Id:              institute.Id.Hex(),
		Name:            institute.Name,
		DisplayableName: institute.DisplayableName,
		MinFloor:        institute.MinFloor,
		MaxFloor:        institute.MaxFloor,
		Url:             institute.Url,
		Latitude:        institute.Latitude,
		Longitude:       institute.Longitude,
		Icon:            utils.IconToIconResponse(icon[0]),
		GPS:             institute.GPS,
	}

	return response, models.ResponseType{Type: 200, Error: nil}
}

func (s *Services) GetAllInstitutes(context context.Context) ([]models.InstituteGet, models.ResponseType) {
	var response []models.InstituteGet

	institutes, res := s.Store.GetAllInstitutes(context)
	if res.Error != nil {
		return response, res
	}
	log.Println(institutes)

	iconNames := []string{}
	for _, institute := range institutes {
		iconNames = append(iconNames, institute.Icon)
	}

	log.Println(iconNames)

	iconsData, iconResp := s.Store.GetInstituteIconsByName(context, iconNames)
	if iconResp.Error != nil {
		log.Println(iconResp)
		return response, res
	}

	if len(iconsData) != len(institutes) {
		log.Printf("IconsData length = %d and InstitutesData length = %d", len(iconsData), len(institutes))
		return response, models.ResponseType{Type: 404, Error: errors.New("For some of the institutes icons not founded")}
	}

	for i, institue := range institutes {
		response = append(response, models.InstituteGet{
			Id:              institue.Id.Hex(),
			Name:            institue.Name,
			DisplayableName: institue.DisplayableName,
			MinFloor:        institue.MinFloor,
			MaxFloor:        institue.MaxFloor,
			Url:             institue.Url,
			Latitude:        institue.Latitude,
			Longitude:       institue.Longitude,
			Icon:            utils.IconToIconResponse(iconsData[i]),
			GPS:             institue.GPS,
		})
	}

	return response, models.ResponseType{Type: 200, Error: nil}
}

func (s *Services) PostInstitute(context context.Context, institute models.InstitutePost) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		return s.Store.InsertInstitute(context, institute), nil
	})
}

func (s *Services) DeleteInstitute(context context.Context, id string) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, nil
		}

		institute, res := s.Store.GetInstitute(ctx, []models.Query{{ParamName: "_id", Type: "ObjectID", ObjectIDValue: objId}})

		if res.Error != nil {
			if res.Error == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified id: " + id)}, nil
			} else {
				return res, nil
			}
		}

		_, res = s.Store.GetFloor(ctx, []models.Query{{ParamName: "institute", Type: "string", StringValue: institute.Name}})
		if res.Error == nil {
			return models.ResponseType{Type: 404, Error: errors.New("can not delete institute with floors")}, nil
		}

		return s.Store.DeleteInstitute(ctx,
			[]models.Query{{ParamName: "_id", Type: "ObjectID", ObjectIDValue: objId}}), nil
	})
}

func (s *Services) UpdateInstitute(context context.Context, body models.InstitutePost, id string) models.ResponseType {
	return s.Store.Transaction(context, func(ctx mongo.SessionContext) (interface{}, error) {
		if body.Name == "" {
			err := errors.New("institute name can not be empty")
			return models.ResponseType{Type: 406, Error: err}, nil
		}

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return models.ResponseType{Type: 500, Error: err}, nil
		}

		filter := []models.Query{{ParamName: "_id", Type: "ObjectID", ObjectIDValue: objId}}

		oldInstitute, res := s.Store.GetInstitute(ctx, filter)
		if res.Error != nil {
			if res.Error == mongo.ErrNoDocuments {
				return models.ResponseType{Type: 404, Error: errors.New("there is no institute with specified id: " + id)}, nil
			} else {
				return res, nil
			}
		}

		if oldInstitute.Name != body.Name || oldInstitute.MaxFloor != body.MaxFloor || oldInstitute.MinFloor != body.MinFloor {
			var floors []models.Floor

			filter = []models.Query{{ParamName: "institute", Type: "string", StringValue: oldInstitute.Name}}
			floors, res := s.Store.GetManyFloors(ctx, filter)
			if res.Error != nil {
				return res, nil
			}

			for _, floor := range floors {
				if body.MinFloor > floor.Floor || floor.Floor > body.MaxFloor {
					return models.ResponseType{Type: 406, Error: errors.New("there is a floor out of new floor bounds")}, nil
				}

				filter = []models.Query{{ParamName: "institute", Type: "string", StringValue: oldInstitute.Name},
					{ParamName: "floor", Type: "int", IntValue: floor.Floor}}

				res = s.Store.UpdateManyGraphs(ctx, filter, struct {
					Institute string `bson:"institute"`
				}{Institute: body.Name})

				if res.Error != nil {
					return res, nil
				}
			}

			filter = []models.Query{{ParamName: "institute", Type: "string", StringValue: oldInstitute.Name}}

			res = s.Store.UpdateManyStairs(ctx, filter, struct {
				Institute string `bson:"institute"`
			}{Institute: body.Name})

			if res.Error != nil {
				return res, nil
			}

			filter = []models.Query{{ParamName: "institute", Type: "string", StringValue: oldInstitute.Name}}

			res = s.Store.UpdateManyFloors(ctx, filter, struct {
				Institute string `bson:"institute"`
			}{Institute: body.Name})

			if res.Error != nil {
				return res, nil
			}
		}

		if oldInstitute.Icon != body.Icon {
			_, res = s.Store.GetInstituteIconsByName(ctx, []string{oldInstitute.Icon})
			if res.Error != nil {
				return res, nil
			}
		}

		filter = []models.Query{{ParamName: "_id", Type: "ObjectID", ObjectIDValue: objId}}
		return s.Store.UpdateInstitute(ctx, filter, body), nil
	})
}
