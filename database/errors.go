package database

import (
	"fmt"

	geojson "github.com/paulmach/go.geojson"
)

// ErrInvalidGeometry ...
type ErrInvalidGeometry struct {
	Type geojson.GeometryType
}

// Error ...
func (err ErrInvalidGeometry) Error() string {
	return fmt.Sprintf("%s is an invalid geometry", err.Type)
}
