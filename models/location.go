package models

// Location ...
type Location struct {
	Lat float64 `json:"lat" validate:"required,gte=-90,lte=90"`
	Lng float64 `json:"lng" validate:"required,gte=-180,lte=180"`
}
