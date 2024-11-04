package models

type Stair struct {
	Id         string   `bson:"_id" json:"id"`
	StairPoint string   `bson:"stairPoint" json:"stairPoint"`
	Institute  string   `json:"institute"`
	Links      []string `json:"links"`
}
