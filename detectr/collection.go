package detectr

import (
	"github.com/tidwall/geoindex"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/rtree"
)

type Collection struct {
	objects int
	index   *geoindex.Index
}

type fence struct {
	object geojson.Object
}

func NewCollection() *Collection {
	c := &Collection{
		index: geoindex.Wrap(&rtree.RTree{}),
	}
	return c
}

func (c *Collection) Insert(g geojson.Object) {
	f := &fence{object: g}

	if !f.object.Empty() {
		rect := f.object.Rect()
		c.index.Insert(
			[2]float64{rect.Min.X, rect.Min.Y},
			[2]float64{rect.Max.X, rect.Max.Y},
			f,
		)
		c.objects++
	}
}

func (c *Collection) Count() int {
	return c.objects
}

func (c *Collection) geoSearch(
	rect geometry.Rect,
	iter func(object geojson.Object) bool,
) bool {
	alive := true
	c.index.Search(
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

func (c *Collection) Intersects(
	obj geojson.Object,
	iter func(object geojson.Object) bool,
) bool {
	return c.geoSearch(obj.Rect(),
		func(f geojson.Object) bool {
			if f.Intersects(obj) {
				return iter(f)
			}
			return true
		},
	)
}
