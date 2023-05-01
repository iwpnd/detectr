package database

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/iwpnd/piper"
	geojson "github.com/paulmach/go.geojson"

	"github.com/tidwall/geoindex"

	"github.com/tidwall/rtree"
)

func toExtent(r [][]float64) []float64 {
	w := r[0][0]
	s := r[0][1]
	e := r[0][0]
	n := r[0][1]

	for _, p := range r {
		if w > p[0] {
			w = p[0]
		}

		if s > p[1] {
			s = p[1]
		}

		if e < p[0] {
			e = p[0]
		}

		if n < p[1] {
			n = p[1]
		}
	}

	return []float64{w, s, e, n}

}

// Database ...
type Database struct {
	tree *geoindex.Index
}

// Fence ...
type Fence struct {
	object geojson.Feature
}

// New to create a new database
func New() *Database {
	db := &Database{
		tree: geoindex.Wrap(&rtree.RTree{}),
	}
	return db
}

// Truncate to create a new database
func (db *Database) Truncate() {
	db.tree = geoindex.Wrap(&rtree.RTree{})
}

// Create to create a new entry into the database
func (db *Database) Create(g *geojson.Feature) error {
	if g.Geometry == nil {
		return &ErrEmptyGeometry{}
	}

	if !g.Geometry.IsPolygon() {
		return &ErrInvalidGeometryType{Type: g.Geometry.Type}
	}

	if g.ID == nil {
		id := uuid.Must(uuid.NewRandom()).String()
		g.ID = id
	}

	f := &Fence{object: *g}

	rect := toExtent(g.Geometry.Polygon[0])
	db.tree.Insert(
		[2]float64{rect[0], rect[1]},
		[2]float64{rect[2], rect[3]},
		f,
	)
	return nil
}

// Delete to delete an entry from the database
func (db *Database) Delete(g *geojson.Feature) {
	rect := toExtent(g.Geometry.Polygon[0])
	db.tree.Delete(
		[2]float64{rect[0], rect[1]},
		[2]float64{rect[2], rect[3]},
		g,
	)
}

// Count to get the current amount of entries in the database
func (db *Database) Count() int {
	return db.tree.Len()
}

func (db *Database) search(
	p []float64,
	iter func(object geojson.Feature) bool,
) bool {
	alive := true
	db.tree.Search(
		[2]float64{p[0], p[1]},
		[2]float64{p[0], p[1]},
		func(_, _ [2]float64, value interface{}) bool {
			item := value.(*Fence)
			alive = iter(item.object)
			return alive
		},
	)
	return alive
}

// Intersects to find entries intersecting the requested point
func (db *Database) Intersects(
	p []float64,
) []geojson.Feature {
	var matches []geojson.Feature

	db.search(p, func(o geojson.Feature) bool {
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
func (db *Database) LoadFromPath(path string) error {
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
