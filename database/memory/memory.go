package memory

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/piper"
	geojson "github.com/paulmach/go.geojson"

	"github.com/tidwall/geoindex"

	"github.com/tidwall/rtree"
)

// Memory ...
type Memory struct {
	tree *geoindex.Index
}

// New to create a new database
func New() *Memory {
	db := &Memory{
		tree: geoindex.Wrap(&rtree.RTree{}),
	}
	return db
}

// Truncate to create a new database
func (db *Memory) Truncate() {
	db.tree = geoindex.Wrap(&rtree.RTree{})
}

// Create to create a new entry into the database
func (db *Memory) Create(g *geojson.Feature) error {
	if g.Geometry == nil {
		return &database.ErrEmptyGeometry{}
	}

	if !g.Geometry.IsPolygon() {
		return &database.ErrInvalidGeometryType{Type: g.Geometry.Type}
	}

	if g.ID == nil {
		id := uuid.Must(uuid.NewRandom()).String()
		g.ID = id
	}

	var or database.OuterRing = g.Geometry.Polygon[0]
	rect := or.ToExtent()
	db.tree.Insert(
		[2]float64{rect[0], rect[1]},
		[2]float64{rect[2], rect[3]},
		g,
	)
	return nil
}

// Delete to delete an entry from the database
func (db *Memory) Delete(g *geojson.Feature) {
	var or database.OuterRing = g.Geometry.Polygon[0]
	rect := or.ToExtent()

	db.tree.Delete(
		[2]float64{rect[0], rect[1]},
		[2]float64{rect[2], rect[3]},
		g,
	)
}

// Count to get the current amount of entries in the database
func (db *Memory) Count() int {
	return db.tree.Len()
}

func (db *Memory) search(
	p []float64,
	iter func(object *geojson.Feature) bool,
) bool {
	alive := true
	db.tree.Search(
		[2]float64{p[0], p[1]},
		[2]float64{p[0], p[1]},
		func(_, _ [2]float64, value interface{}) bool {
			item := value.(*geojson.Feature)
			alive = iter(item)
			return alive
		},
	)
	return alive
}

// Intersects to find entries intersecting the requested point
func (db *Memory) Intersects(p []float64) []*geojson.Feature {
	var matches []*geojson.Feature

	db.search(p, func(o *geojson.Feature) bool {
		if o.Geometry.IsPolygon() {
			if piper.Pip(p, o.Geometry.Polygon) {
				matches = append(matches, o)
			}
			return true
		}
		return true
	})

	return matches
}

// LoadFromPath to load a FeatureCollection from file
func (db *Memory) LoadFromPath(path string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	fc, err := geojson.UnmarshalFeatureCollection(file)
	if err != nil {
		return err
	}

	for _, f := range fc.Features {
		err := db.Create(f)
		if err != nil {
			fmt.Print("skipping geometry: ", err.Error())
		}
	}

	return nil
}
