package models

import geojson "github.com/paulmach/go.geojson"

// Response ...
type Response struct {
	Data interface{} `json:"data"`
}

// LocationResponsePayload ...
type LocationResponsePayload struct {
	Elapsed string            `json:"elapsed"`
	Request Location          `json:"request"`
	Matches []geojson.Feature `json:"matches"`
}

// LocationResponse ...
type LocationResponse struct {
	Data LocationResponsePayload `json:"data"`
}
