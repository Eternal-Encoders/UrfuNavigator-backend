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

func GraphPointToGraphPointPut(gp models.GraphPoint) models.GraphPointPut {
	return models.GraphPointPut{
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

func FloorToFloorPut(floor models.Floor) models.FloorPut {
	return models.FloorPut{
		Institute: floor.Institute,
		Floor:     floor.Floor,
		Width:     floor.Height,
		Height:    floor.Height,
		Audiences: floor.Audiences,
		Service:   floor.Service,
		Graph:     floor.Graph,
		Forces:    floor.Forces,
	}
}
