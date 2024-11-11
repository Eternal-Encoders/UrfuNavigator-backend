package models

type InstituteReadDBResponse struct {
	Response  Institute
	Error     error
	ErrorType int
}
