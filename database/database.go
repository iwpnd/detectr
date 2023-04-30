package database

import (
	"os"

	"github.com/google/uuid"
	"github.com/tidwall/geoindex"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/rtree"
)

type Database struct {
	tree *geoindex.Index
}

type fence struct {
	id     string
	object geojson.Object
}

func New() *Database {
	db := &Database{
		tree: geoindex.Wrap(&rtree.RTree{}),
	}
	return db
}

func (db *Database) Truncate() {
	db.tree = geoindex.Wrap(&rtree.RTree{})
}

func (db *Database) Create(g geojson.Object) {
	id := uuid.Must(uuid.NewRandom()).String()
	f := &fence{object: g, id: id}

	if !f.object.Empty() {
		rect := f.object.Rect()
		db.tree.Insert(
			[2]float64{rect.Min.X, rect.Min.Y},
			[2]float64{rect.Max.X, rect.Max.Y},
			f,
		)
	}
}

func (db *Database) Delete(g geojson.Object) {
	rect := g.Rect()
	db.tree.Delete(
		[2]float64{rect.Min.X, rect.Min.Y},
		[2]float64{rect.Max.X, rect.Max.Y},
		g,
	)
}

func (db *Database) Count() int {
	return db.tree.Len()
}

func (db *Database) search(
	rect geometry.Rect,
	iter func(object geojson.Object) bool,
) bool {
	alive := true
	db.tree.Search(
		[2]float64{rect.Min.X, rect.Min.Y},
		[2]float64{rect.Max.X, rect.Max.Y},
		func(_, _ [2]float64, value interface{}) bool {
			item := value.(*fence)
			alive = iter(item.object)
			return alive
		},
	)
	return alive
}

func (db *Database) Intersects(
	obj geojson.Object,
) []geojson.Object {
	var matches []geojson.Object

	db.search(obj.Rect(), func(o geojson.Object) bool {
		if obj.Intersects(o) {
			matches = append(matches, o)
		}
		return true
	})

	return matches
}

func (db *Database) LoadFromPath(path string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	fc, err := geojson.Parse(string(file), nil)
	if err != nil {
		return err
	}

	fc.ForEach(func(o geojson.Object) bool {
		if o.Empty() {
			return true
		}

		db.Create(o)
		return true
	})

	return nil
}
