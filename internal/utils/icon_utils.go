package utils

import "UrfuNavigator-backend/internal/models"

func IconToIconResponse(icon models.InstituteIcon) models.InstituteIconResponse {
	return models.InstituteIconResponse{
		Id:  icon.Id.Hex(),
		Url: icon.Url,
		Alt: icon.Alt,
	}
}
