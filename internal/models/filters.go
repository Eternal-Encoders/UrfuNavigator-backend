package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Filter struct {
	StringFilters []StringFilter
	IntFilters    []IntFilter
}

type StringFilter struct {
	ParamName string
	Value     string
}

type IntFilter struct {
	ParamName string
	Value     int
}

type Query struct {
	ParamName     string
	Type          string
	IntValue      int
	FloatValue    float32
	StringValue   string
	ObjectIDValue primitive.ObjectID
}
