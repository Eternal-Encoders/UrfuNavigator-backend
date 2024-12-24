package models

type User struct {
	Email    string `bson:"email" json:"email"`
	Password string `json:"password"`
}

type UserDB struct {
	Email string `bson:"email" json:"email"`
	Hash  string `json:"hash"`
	// Salt  string `json:"salt"`
}
