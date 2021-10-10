package fences

import (
	"github.com/tidwall/geoindex"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/rtree"
	"io/ioutil"
)

type Fences struct {
	objects int
	tree    *geoindex.Index
}

type fence struct {
	object geojson.Object
}

var database = New()

func New() *Fences {
	fences := &Fences{
		tree: geoindex.Wrap(&rtree.RTree{}),
	}
	return fences
}

func Get() *Fences {
	return database
}

func (fences *Fences) Create(g geojson.Object) {
	f := &fence{object: g}

	if !f.object.Empty() {
		rect := f.object.Rect()
		fences.tree.Insert(
			[2]float64{rect.Min.X, rect.Min.Y},
			[2]float64{rect.Max.X, rect.Max.Y},
			f,
		)
		fences.objects++
	}
}

func (fences *Fences) Count() int {
	return fences.objects
}

func (fences *Fences) search(
	rect geometry.Rect,
	iter func(object geojson.Object) bool,
) bool {
	alive := true
	fences.tree.Search(
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

func (fences *Fences) intersects(
	obj geojson.Object,
	iter func(object geojson.Object) bool,
) bool {
	return fences.search(obj.Rect(),
		func(f geojson.Object) bool {
			if f.Intersects(obj) {
				return iter(f)
			}
			return true
		},
	)
}

func (fences *Fences) Intersects(
	obj geojson.Object,
) []geojson.Object {
	var matches []geojson.Object

	fences.intersects(obj, func(o geojson.Object) bool {
		matches = append(matches, o)
		return true
	})

	return matches
}

func (fences *Fences) LoadFromPath(path string) error {
	file, err := ioutil.ReadFile(path)
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

		fences.Create(o)
		fences.objects++
		return true
	})

	return nil
}
