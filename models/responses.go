package models

import "github.com/tidwall/geojson"

type Response struct {
	Data interface{} `json:"data"`
}

type LocationResponsePayload struct {
	Elapsed string           `json:"elapsed"`
	Request Location         `json:"request"`
	Matches []geojson.Object `json:"matches"`
}

type LocationResponse struct {
	Data LocationResponsePayload `json:"data"`
}
