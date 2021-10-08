package detectr

import (
	"github.com/buckhx/diglet/geo"
)

type Collection struct {
	index *geo.Rtree
}

func NewCollection() *Collection {
	c := &Collection{
		index: geo.NewRtree(),
	}
	return c
}

func (c *Collection) Insert(f *geo.Feature) {
	for _, p := range f.Geometry {
		if len(p.Coordinates) > 1 {
			c.index.Insert(p, f)
		}
	}
}

func (c *Collection) Contains(coord geo.Coordinate) (matches []*geo.Feature) {
	nodes := c.index.Contains(coord)
	for _, n := range nodes {
		feature := n.Feature()
		if feature.Contains(coord) {
			matches = append(matches, feature)
		}
	}

	return
}
