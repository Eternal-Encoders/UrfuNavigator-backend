package models

import "errors"

type ResponseType struct {
	Type  int
	Error error
}

var ErrNoDocuments = errors.New("mongo: no documents in result")
