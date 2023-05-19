package database

import geojson "github.com/paulmach/go.geojson"

type Searcher interface {
	Intersects(p []float64) []geojson.Feature
}

type Creater interface {
	Create(*geojson.Feature) error
}

type Deleter interface {
	Delete(*geojson.Feature)
	Truncate()
}

type Datastore interface {
	Searcher
	Creater
	Deleter
}

func ToExtent(r [][]float64) []float64 {
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
