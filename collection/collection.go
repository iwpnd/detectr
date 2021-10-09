package collection

import (
	"fmt"
	"github.com/tidwall/geoindex"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/rtree"
	"io/ioutil"
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

func (c *Collection) intersects(
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

func (c *Collection) Intersects(
	obj geojson.Object,
) []geojson.Object {
	var items []geojson.Object

	c.intersects(obj, func(o geojson.Object) bool {
		items = append(items, o)
		return true
	})

	return items
}

func (c *Collection) LoadFromPath(path string) error {
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

		fmt.Println(o)

		c.Insert(o)
		return true
	})

	return nil
}