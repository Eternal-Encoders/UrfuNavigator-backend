package models

type WeekTime struct {
	From    string `json:"from"`
	To      string `json:"to"`
	IsDayOf *bool  `json:"isDayOf"`
}

type GraphPoint struct {
	Id          string     `bson:"_id" json:"id"`
	X           float64    `json:"x"`
	Y           float64    `json:"y"`
	Links       []string   `json:"links"`
	Types       []string   `json:"types"`
	Names       []string   `json:"names"`
	Floor       int        `json:"floor"`
	Institute   string     `json:"institute"`
	Time        []WeekTime `json:"time"`
	Description string     `json:"description"`
	Info        string     `json:"info"`
	MenuId      *string    `json:"menuId"`
	IsPassFree  *bool      `json:"isPassFree"`
	StairId     *string    `json:"stairId"`
}
