package models

type Location struct {
	Lat float64 `validate:"required,gte=-90,lte=90"`
	Lng float64 `validate:"required,gte=-180,lte=180"`
}
