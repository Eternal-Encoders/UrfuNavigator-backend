package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type InstituteIconPost struct {
	Id  string `bson:"id" json:"id"`
	Url string `bson:"filename" json:"url"`
	Alt string `json:"alt"`
}

type InstituteIcon struct {
	Id  primitive.ObjectID `bson:"_id" json:"id"`
	Url string             `bson:"filename" json:"url"`
	Alt string             `json:"alt"`
}

type InstituteIconGet struct {
	Url string `bson:"filename" json:"url"`
	Alt string `json:"alt"`
}
