package database

import (
	"fmt"

	geojson "github.com/paulmach/go.geojson"
)

// ErrInvalidGeometry ...
type ErrInvalidGeometryType struct {
	Type geojson.GeometryType
}

// Error ...
func (err ErrInvalidGeometryType) Error() string {
	return fmt.Sprintf("unsupported geometry type: %s", err.Type)
}

// ErrEmptyGeometry ...
type ErrEmptyGeometry struct{}

// Error ...
func (err ErrEmptyGeometry) Error() string {
	return "empty geometry"
}
