package geoindex

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"testing"

	"github.com/tidwall/cities"
	"github.com/tidwall/geoindex/algo"
	"github.com/tidwall/lotsa"
)

// Tests is a set of tests for running against an object that conforms to
// geoindex.Interface.
// These functions are intended to be included in a `_test.go` file. A complete
// test file might look something like:
//
// 		package myrtree
//
// 		import (
// 			"math/rand"
// 			"testing"
// 			"time"
//
// 			"github.com/tidwall/geoindex"
// 		)
//
// 		func init() {
// 			seed := time.Now().UnixNano()
// 			println("seed:", seed)
// 			rand.Seed(seed)
// 		}
//
// 		func TestGeoIndex(t *testing.T) {
// 			t.Run("BenchInsert", func(t *testing.T) {
// 				geoindex.Tests.TestBenchInsert(t, &RTree{}, 100000)
// 			})
// 			t.Run("RandomRects", func(t *testing.T) {
// 				geoindex.Tests.TestRandomRects(t, &RTree{}, 10000)
// 			})
// 			t.Run("RandomPoints", func(t *testing.T) {
// 				geoindex.Tests.TestRandomPoints(t, &RTree{}, 10000)
// 			})
// 			t.Run("ZeroPoints", func(t *testing.T) {
// 				geoindex.Tests.TestZeroPoints(t, &RTree{})
// 			})
// 			t.Run("CitiesSVG", func(t *testing.T) {
// 				geoindex.Tests.TestCitiesSVG(t, &RTree{})
// 			})
// 		}
//
// 		func BenchmarkRandomInsert(b *testing.B) {
// 			geoindex.Tests.BenchmarkRandomInsert(b, &RTree{})
// 		}
//
var Tests = struct {
	TestBenchVarious      func(t *testing.T, tr Interface, numPointOrRects int)
	TestRandomPoints      func(t *testing.T, tr Interface, numPoints int)
	TestRandomRects       func(t *testing.T, tr Interface, numRects int)
	TestCitiesSVG         func(t *testing.T, tr Interface)
	TestZeroPoints        func(t *testing.T, tr Interface)
	BenchmarkRandomInsert func(b *testing.B, tr Interface)
}{
	benchVarious,
	func(t *testing.T, tr Interface, numRects int) {
		testBoxesVarious(t, tr, randBoxes(numRects), "boxes")
	},
	func(t *testing.T, tr Interface, numPoints int) {
		testBoxesVarious(t, tr, randPoints(numPoints), "points")
	},
	testCitiesSVG,
	testZeroPoints,
	benchmarkRandomInsert,
}

type rect struct {
	min, max [2]float64
}

// kind = 'r','p','m' for rect,point,mixed
func randRect(kind byte) (r rect) {
	r.min[0] = rand.Float64()*360 - 180
	r.min[1] = rand.Float64()*180 - 90
	r.max = r.min
	return randRectOffset(r, kind)
}

func randRectOffset(r rect, kind byte) rect {
	rsize := 0.01 // size of rectangle in degrees
	pr := r
	for {
		r.min[0] = (pr.max[0]+pr.min[0])/2 + rand.Float64()*rsize - rsize/2
		r.min[1] = (pr.max[1]+pr.min[1])/2 + rand.Float64()*rsize - rsize/2
		r.max = r.min
		if kind == 'r' || (kind == 'm' && rand.Int()%2 == 0) {
			// rect
			r.max[0] = r.min[0] + rand.Float64()*rsize
			r.max[1] = r.min[1] + rand.Float64()*rsize
		} else {
			// point
			r.max = r.min
		}
		if r.min[0] < -180 || r.min[1] < -90 ||
			r.max[0] > 180 || r.max[1] > 90 {
			continue
		}
		return r
	}
}

type mixedTree interface {
	IsMixedTree() bool
}

func benchVarious(t *testing.T, tr Interface, numPointOrRects int) {
	if v, ok := tr.(mixedTree); ok && v.IsMixedTree() {
		println("== points ==")
		benchVariousKind(t, tr, numPointOrRects, 'p')
		println("== rects ==")
		benchVariousKind(t, tr, numPointOrRects, 'r')
		println("== mixed (50/50) ==")
		benchVariousKind(t, tr, numPointOrRects, 'm')
	} else {
		benchVariousKind(t, tr, numPointOrRects, 'm')
	}
}

func benchVariousKind(t *testing.T, tr Interface, numPointOrRects int,
	kind byte,
) {
	N := numPointOrRects
	rects := make([]rect, N)
	for i := 0; i < N; i++ {
		rects[i] = randRect(kind)
	}
	rectsReplace := make([]rect, N)
	for i := 0; i < N; i++ {
		rectsReplace[i] = randRectOffset(rects[i], kind)
	}
	lotsa.Output = os.Stdout
	fmt.Printf("insert:  ")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Insert(rects[i].min, rects[i].max, i)
	})
	fmt.Printf("search:  ")
	var count int
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Search(rects[i].min, rects[i].max,
			func(min, max [2]float64, value interface{}) bool {
				if value.(int) == i {
					count++
					return false
				}
				return true
			},
		)
	})
	if count != N {
		t.Fatalf("expected %d, got %d", N, count)
	}
	fmt.Printf("replace: ")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Replace(
			rects[i].min, rects[i].max, i,
			rectsReplace[i].min, rectsReplace[i].max, i,
		)
	})
	if tr.Len() != N {
		t.Fatalf("expected %d, got %d", N, tr.Len())
	}

	fmt.Printf("delete:  ")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Delete(rectsReplace[i].min, rectsReplace[i].max, i)
	})
	if tr.Len() != 0 {
		t.Fatalf("expected %d, got %d", 0, tr.Len())
	}
}

func testBoxesVarious(t *testing.T, tr Interface, boxes []tBox, label string) {
	N := len(boxes)

	/////////////////////////////////////////
	// insert
	/////////////////////////////////////////
	for i := 0; i < N; i++ {
		tr.Insert(boxes[i].min, boxes[i].max, boxes[i])
	}
	if tr.Len() != N {
		t.Fatalf("expected %d, got %d", N, tr.Len())
	}

	/////////////////////////////////////////
	// scan all items and count one-by-one
	/////////////////////////////////////////
	var count int
	tr.Scan(func(min, max [2]float64, value interface{}) bool {
		count++
		return true
	})
	if count != N {
		t.Fatalf("expected %d, got %d", N, count)
	}

	/////////////////////////////////////////
	// check every point for correctness
	/////////////////////////////////////////
	var tboxes1 []tBox
	tr.Scan(func(min, max [2]float64, value interface{}) bool {
		tboxes1 = append(tboxes1, value.(tBox))
		return true
	})
	tboxes2 := make([]tBox, len(boxes))
	copy(tboxes2, boxes)
	sortBoxes(tboxes1)
	sortBoxes(tboxes2)
	for i := 0; i < len(tboxes1); i++ {
		if tboxes1[i] != tboxes2[i] {
			t.Fatalf("expected '%v', got '%v'", tboxes2[i], tboxes1[i])
		}
	}

	/////////////////////////////////////////
	// search for each item one-by-one
	/////////////////////////////////////////
	for i := 0; i < N; i++ {
		var found bool
		tr.Search(boxes[i].min, boxes[i].max,
			func(min, max [2]float64, value interface{}) bool {
				if value == boxes[i] {
					found = true
					return false
				}
				return true
			})
		if !found {
			t.Fatalf("did not find item %d", i)
		}
	}

	centerMin, centerMax := [2]float64{-18, -9}, [2]float64{18, 9}

	/////////////////////////////////////////
	// search for 10% of the items
	/////////////////////////////////////////
	for i := 0; i < N/5; i++ {
		var count int
		tr.Search(centerMin, centerMax,
			func(min, max [2]float64, value interface{}) bool {
				count++
				return true
			},
		)
	}

	/////////////////////////////////////////
	// delete every other item
	/////////////////////////////////////////
	for i := 0; i < N/2; i++ {
		j := i * 2
		tr.Delete(boxes[j].min, boxes[j].max, boxes[j])
	}

	/////////////////////////////////////////
	// count all items. should be half of N
	/////////////////////////////////////////
	count = 0
	tr.Scan(func(min, max [2]float64, value interface{}) bool {
		count++
		return true
	})
	if count != N/2 {
		t.Fatalf("expected %d, got %d", N/2, count)
	}

	///////////////////////////////////////////////////
	// reinsert every other item, but in random order
	///////////////////////////////////////////////////
	var ij []int
	for i := 0; i < N/2; i++ {
		j := i * 2
		ij = append(ij, j)
	}
	rand.Shuffle(len(ij), func(i, j int) {
		ij[i], ij[j] = ij[j], ij[i]
	})
	for i := 0; i < N/2; i++ {
		j := ij[i]
		tr.Insert(boxes[j].min, boxes[j].max, boxes[j])
	}

	//////////////////////////////////////////////////////
	// replace each item with an item that is very close
	//////////////////////////////////////////////////////
	var nboxes = make([]tBox, N)
	for i := 0; i < N; i++ {
		for j := 0; j < len(boxes[i].min); j++ {
			nboxes[i].min[j] = boxes[i].min[j] + (rand.Float64() - 0.5)
			if boxes[i].min == boxes[i].max {
				nboxes[i].max[j] = nboxes[i].min[j]
			} else {
				nboxes[i].max[j] = boxes[i].max[j] + (rand.Float64() - 0.5)
			}
		}

	}
	for i := 0; i < N; i++ {
		tr.Insert(nboxes[i].min, nboxes[i].max, nboxes[i])
		tr.Delete(boxes[i].min, boxes[i].max, boxes[i])
	}
	if tr.Len() != N {
		t.Fatalf("expected %d, got %d", N, tr.Len())
	}

	/////////////////////////////////////////
	// check every point for correctness
	/////////////////////////////////////////
	tboxes1 = nil
	tr.Scan(func(min, max [2]float64, value interface{}) bool {
		tboxes1 = append(tboxes1, value.(tBox))
		return true
	})
	tboxes2 = make([]tBox, len(nboxes))
	copy(tboxes2, nboxes)
	sortBoxes(tboxes1)
	sortBoxes(tboxes2)
	for i := 0; i < len(tboxes1); i++ {
		if tboxes1[i] != tboxes2[i] {
			t.Fatalf("expected '%v', got '%v'", tboxes2[i], tboxes1[i])
		}
	}

	/////////////////////////////////////////
	// search for 10% of the items
	/////////////////////////////////////////
	for i := 0; i < N/5; i++ {
		var count int
		tr.Search(centerMin, centerMax,
			func(min, max [2]float64, value interface{}) bool {
				count++
				return true
			},
		)
	}

	var boxes3 []tBox
	Wrap(tr).Nearby(
		algo.Box(centerMin, centerMax, false, nil),
		func(min, max [2]float64, value interface{}, dist float64) bool {
			boxes3 = append(boxes3, value.(tBox))
			return true
		},
	)

	if len(boxes3) != len(nboxes) {
		t.Fatalf("expected %d, got %d", len(nboxes), len(boxes3))
	}
	if len(boxes3) != tr.Len() {
		t.Fatalf("expected %d, got %d", tr.Len(), len(boxes3))
	}

	var ldist float64
	for i, box := range boxes3 {
		dist := testBoxDist(box.min, box.max, centerMin, centerMax)
		if i > 0 && dist < ldist {
			t.Fatalf("out of order")
		}
		ldist = dist
	}
}

func sortBoxes(boxes []tBox) {
	sort.Slice(boxes, func(i, j int) bool {
		for k := 0; k < len(boxes[i].min); k++ {
			if boxes[i].min[k] < boxes[j].min[k] {
				return true
			}
			if boxes[i].min[k] > boxes[j].min[k] {
				return false
			}
			if boxes[i].max[k] < boxes[j].max[k] {
				return true
			}
			if boxes[i].max[k] > boxes[j].max[k] {
				return false
			}
		}
		return i < j
	})
}

func testBoxDist(amin, amax, bmin, bmax [2]float64) float64 {
	var dist float64
	for i := 0; i < len(amin); i++ {
		var min, max float64
		if amin[i] > bmin[i] {
			min = amin[i]
		} else {
			min = bmin[i]
		}
		if amax[i] < bmax[i] {
			max = amax[i]
		} else {
			max = bmax[i]
		}
		squared := min - max
		if squared > 0 {
			dist += squared * squared
		}
	}
	return dist
}

func randPoints(N int) []tBox {
	boxes := make([]tBox, N)
	for i := 0; i < N; i++ {
		boxes[i].min[0] = rand.Float64()*360 - 180
		boxes[i].min[1] = rand.Float64()*180 - 90
		boxes[i].max = boxes[i].min
	}
	return boxes
}

func testZeroPoints(t *testing.T, tr Interface) {
	N := 10000
	var pt [2]float64
	for i := 0; i < N; i++ {
		tr.Insert(pt, pt, i)
	}
}

func benchmarkRandomInsert(b *testing.B, tr Interface) {
	boxes := randBoxes(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Insert(boxes[i].min, boxes[i].max, i)
	}
}

func testCitiesSVG(t *testing.T, tr Interface) {
	index := Wrap(tr)
	for _, i := range rand.Perm(len(cities.Cities)) {
		city := cities.Cities[i]
		p := [2]float64{city.Longitude, city.Latitude}
		index.Insert(p, p, &city)
	}
	svg := index.SVG()
	if err := ioutil.WriteFile("cities.svg", []byte(svg), 0600); err != nil {
		t.Fatal(err)
	}
}

type tBox struct {
	min [2]float64
	max [2]float64
}

func randBoxes(N int) []tBox {
	boxes := make([]tBox, N)
	for i := 0; i < N; i++ {
		boxes[i].min[0] = rand.Float64()*360 - 180
		boxes[i].min[1] = rand.Float64()*180 - 90
		boxes[i].max[0] = boxes[i].min[0] + rand.Float64()
		boxes[i].max[1] = boxes[i].min[1] + rand.Float64()
		if boxes[i].max[0] > 180 || boxes[i].max[1] > 90 {
			i--
		}
	}
	return boxes
}
