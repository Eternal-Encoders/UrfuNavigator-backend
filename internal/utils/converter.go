package utils

import "UrfuNavigator-backend/internal/models"

func IconToIconResponse(icon models.InstituteIcon) models.InstituteIconPost {
	return models.InstituteIconPost{
		Id:  icon.Id.Hex(),
		Url: icon.Url,
		Alt: icon.Alt,
	}
}

func GraphPointPutToGraphPoint(gp models.GraphPointPut, id string) models.GraphPoint {
	return models.GraphPoint{
		Id:          id,
		X:           gp.X,
		Y:           gp.Y,
		Links:       gp.Links,
		Types:       gp.Types,
		Names:       gp.Names,
		Floor:       gp.Floor,
		Institute:   gp.Institute,
		Time:        gp.Time,
		Description: gp.Description,
		Info:        gp.Info,
		MenuId:      gp.MenuId,
		IsPassFree:  gp.IsPassFree,
		StairId:     gp.StairId,
	}
}
