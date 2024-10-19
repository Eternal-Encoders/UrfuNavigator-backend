package models

type Floor struct {
	Institute string     `json:"institute"`
	Floor     int        `json:"floor"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Audiences []Audience `json:"audiences"`
	Service   []Service  `json:"service"`
	Graph     []string   `json:"graph"`
}

type FloorFromFile struct {
	Institute string       `json:"institute"`
	Floor     int          `json:"floor"`
	Width     int          `json:"width"`
	Height    int          `json:"height"`
	Audiences []Audience   `json:"audiences"`
	Service   []Service    `json:"service"`
	Graph     []GraphPoint `json:"graph"`
}

type Audience struct {
	Id       string          `json:"id"`
	X        float64         `json:"x"`
	Y        float64         `json:"y"`
	Width    float64         `json:"width"`
	Height   float64         `json:"height"`
	Fill     string          `json:"fill"`
	Stroke   string          `json:"stroke"`
	PointId  string          `json:"pointId"`
	Children []AudienceChild `json:"children"`
	Doors    []Door          `json:"doors"`
}

type AudienceChild struct {
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
	Stroke *string `json:"stroke"`
	Fill   *string `json:"fill"`
}
