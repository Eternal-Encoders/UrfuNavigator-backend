package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Institute struct {
	Id              primitive.ObjectID `bson:"_id" json:"id"`
	Name            string             `bson:"name" json:"name"`
	DisplayableName string             `bson:"displayableName" json:"displayableName"`
	MinFloor        int                `bson:"minFloor" json:"minFloor"`
	MaxFloor        int                `bson:"maxFloor" json:"maxFloor"`
	Url             string             `bson:"url" json:"url"`
	Latitude        float64            `json:"latitude"`
	Longitude       float64            `json:"longitude"`
	Icon            string             `json:"icon"`
}

type InstituteGet struct {
	Id              string            `json:"id"`
	Name            string            `json:"name"`
	DisplayableName string            `bson:"displayableName" json:"displayableName"`
	MinFloor        int               `bson:"minFloor" json:"minFloor"`
	MaxFloor        int               `bson:"maxFloor" json:"maxFloor"`
	Url             string            `json:"url"`
	Latitude        float64           `json:"latitude"`
	Longitude       float64           `json:"longitude"`
	Icon            InstituteIconPost `json:"icon"`
}

type InstitutePost struct {
	Name            string  `json:"name"`
	DisplayableName string  `bson:"displayableName" json:"displayableName"`
	MinFloor        int     `bson:"minFloor" json:"minFloor"`
	MaxFloor        int     `bson:"maxFloor" json:"maxFloor"`
	Url             string  `json:"url"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Icon            string  `json:"icon"`
}

type InstituteNullable struct {
	Name            *string  `json:"name"`
	DisplayableName *string  `bson:"displayableName" json:"displayableName"`
	MinFloor        *int     `bson:"minFloor" json:"minFloor"`
	MaxFloor        *int     `bson:"maxFloor" json:"maxFloor"`
	Url             *string  `json:"url"`
	Latitude        *float64 `json:"latitude"`
	Longitude       *float64 `json:"longitude"`
	Icon            *string  `json:"icon"`
}
