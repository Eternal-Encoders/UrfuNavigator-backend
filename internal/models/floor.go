package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Floor struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Institute string             `json:"institute"`
	Floor     int                `json:"floor"`
	Width     int                `json:"width"`
	Height    int                `json:"height"`
	Audiences []*Auditorium      `json:"audiences"`
	Service   []*Service         `json:"service"`
	Graph     []string           `json:"graph"`
	Forces    []*Forces          `bson:"forces" json:"forces"`
}

type FloorResponse struct {
	Id        string        `json:"id"`
	Institute string        `json:"institute"`
	Floor     int           `json:"floor"`
	Width     int           `json:"width"`
	Height    int           `json:"height"`
	Audiences []*Auditorium `json:"audiences"`
	Service   []*Service    `json:"service"`
	Graph     []string      `json:"graph"`
	Forces    []*Forces     `bson:"forces" json:"forces"`
}

type FloorRequest struct {
	Institute string        `json:"institute"`
	Floor     int           `json:"floor"`
	Width     int           `json:"width"`
	Height    int           `json:"height"`
	Audiences []*Auditorium `json:"audiences"`
	Service   []*Service    `json:"service"`
	Graph     []string      `json:"graph"`
	Forces    []*Forces     `bson:"forces" json:"forces"`
}

type FloorPut struct {
	Institute string        `json:"institute"`
	Floor     int           `json:"floor"`
	Width     int           `json:"width"`
	Height    int           `json:"height"`
	Audiences []*Auditorium `json:"audiences"`
	Service   []*Service    `json:"service"`
	Graph     []string      `json:"graph"`
	Forces    []*Forces     `bson:"forces" json:"forces"`
}

type FloorFromFile struct {
	Institute string                 `json:"institute"`
	Floor     int                    `json:"floor"`
	Width     int                    `json:"width"`
	Height    int                    `json:"height"`
	Audiences map[string]*Auditorium `json:"audiences"`
	Service   []*Service             `json:"service"`
	Graph     map[string]*GraphPoint `json:"graph"`
	Forces    []*Forces              `bson:"forces" json:"forces"`
}

type Auditorium struct {
	Id       string             `json:"id"`
	X        float64            `json:"x"`
	Y        float64            `json:"y"`
	Width    float64            `json:"width"`
	Height   float64            `json:"height"`
	Fill     string             `json:"fill"`
	Stroke   string             `json:"stroke"`
	PointId  string             `json:"pointId"`
	Children []*AuditoriumChild `json:"children"`
	Doors    []*Door            `json:"doors"`
}

type AuditoriumChild struct {
	Type       string  `json:"type"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Identifier string  `json:"identifier"`
	AlignX     string  `json:"alignX"`
	AlignY     string  `json:"alignY"`
}

type Door struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
	Fill   string  `json:"fill"`
}

type Service struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Data   string  `json:"data"`
	Stroke string  `json:"stroke"`
	Fill   string  `json:"fill"`
}

type Forces struct {
	Point Coordinates `bson:"point" json:"point"`
	Force Coordinates `bson:"force" json:"force"`
}

type Coordinates struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}
