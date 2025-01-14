package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Email    string `bson:"email" json:"email"`
	Password string `json:"password"`
}

// type UserDB struct {
// 	Email string `bson:"email" json:"email"`
// 	Hash  string `json:"hash"`
// 	// Salt  string `json:"salt"`
// }

type UserDB struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Email string             `bson:"email" json:"email"`
	Hash  string             `json:"hash"`
}

type UserCreate struct {
	Email string `bson:"email" json:"email"`
	Hash  string `json:"hash"`
}

type UserDTO struct {
	Id    string `bson:"_id" json:"id"`
	Email string `bson:"email" json:"email"`
}
