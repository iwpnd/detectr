package models

import "github.com/tidwall/geojson"

type Response struct {
	Data interface{} `json:"data"`
}

type GeofenceResponse struct {
	Elapsed string           `json:"elapsed"`
	Request Location         `json:"request"`
	Matches []geojson.Object `json:"matches"`
}
