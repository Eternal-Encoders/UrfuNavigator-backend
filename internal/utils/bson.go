package utils

import (
	"UrfuNavigator-backend/internal/models"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// func CreateBSONFilter(filter models.Filter) (bson.D, error) {
// 	var bsonFilter bson.D
// 	for _, v := range filter.IntFilters {
// 		bsonFilter = append(bsonFilter, bson.E{v.ParamName, v.Value})
// 	}

// 	for _, v := range filter.StringFilters {
// 		if v.ParamName == "_id" {
// 			objId, err := primitive.ObjectIDFromHex(v.Value)
// 			if err != nil {
// 				return bsonFilter, err
// 			}

// 			bsonFilter = append(bsonFilter, bson.E{v.ParamName, objId})
// 		} else {
// 			bsonFilter = append(bsonFilter, bson.E{v.ParamName, v.Value})
// 		}
// 	}

// 	return bsonFilter, nil
// }

func CreateBSONFilter(filter []models.Query) (bson.D, error) {
	var bsonFilter bson.D
	for _, v := range filter {
		switch v.Type {
		case "string":
			bsonFilter = append(bsonFilter, bson.E{v.ParamName, v.StringValue})
		case "float32":
			bsonFilter = append(bsonFilter, bson.E{v.ParamName, v.FloatValue})
		case "int":
			bsonFilter = append(bsonFilter, bson.E{v.ParamName, v.IntValue})
		case "ObjectID":
			bsonFilter = append(bsonFilter, bson.E{v.ParamName, v.ObjectIDValue})
		default:
			return bsonFilter, errors.New("wrong filter type")
		}
	}

	return bsonFilter, nil
}

func ToBson(e interface{}) (data bson.M) {
	var tagValue string

	data = bson.M{}
	element := reflect.ValueOf(e).Elem()

	for i := 0; i < element.NumField(); i += 1 {
		typeField := element.Type().Field(i)
		tag := typeField.Tag

		tagValue = tag.Get("bson")

		if tagValue == "-" {
			continue
		}

		switch element.Field(i).Kind() {

		case reflect.String:
			value := element.Field(i).String()
			data[tagValue] = value

		case reflect.Float32:
			value := element.Field(i).Float()
			data[tagValue] = value

		case reflect.Int:
			value := element.Field(i).Int()
			data[tagValue] = value

		case reflect.Slice:
			valueLen := element.Field(i).Len()
			value := element.Field(i).Slice(0, valueLen)
			data[tagValue] = value
		}
	}

	return data
}
