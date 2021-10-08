package geoindex

import (
	"fmt"

	"github.com/tidwall/geoindex/child"
)

// Interface is a tree-like structure that contains geospatial data
type Interface interface {
	// Insert an item into the structure
	Insert(min, max [2]float64, data interface{})
	// Delete an item from the structure
	Delete(min, max [2]float64, data interface{})
	// Replace an item in the structure. This is effectively just a Delete
	// followed by an Insert. But for some structures it may be possible to
	// optimize the operation to avoid multiple passes
	Replace(
		oldMin, oldMax [2]float64, oldData interface{},
		newMin, newMax [2]float64, newData interface{},
	)
	// Search the structure for items that intersects the rect param
	Search(
		min, max [2]float64,
		iter func(min, max [2]float64, data interface{}) bool,
	)
	// Scan iterates through all data in tree in no specified order.
	Scan(iter func(min, max [2]float64, data interface{}) bool)
	// Len returns the number of items in tree
	Len() int
	// Bounds returns the minimum bounding box
	Bounds() (min, max [2]float64)
	// Children returns all children for parent node. If parent node is nil
	// then the root nodes should be returned.
	// The reuse buffer is an empty length slice that can optionally be used
	// to avoid extra allocations.
	Children(parent interface{}, reuse []child.Child) (children []child.Child)
}

// Index is a wrapper around Interface that provides extra features like a
// Nearby (kNN) function.
// This can be created like such:
//   var tree = &rtree.RTree{}
//   var index = index.Index{tree}
// Now you can use `index` just like tree but with the extra features.
type Index struct {
	tree Interface
}

// Wrap a tree-like geospatial interface.
func Wrap(tree Interface) *Index {
	return &Index{tree}
}

// Insert an item into the index
func (index *Index) Insert(min, max [2]float64, data interface{}) {
	index.tree.Insert(min, max, data)
}

// Search the index for items that intersects the rect param
func (index *Index) Search(
	min, max [2]float64,
	iter func(min, max [2]float64, data interface{}) bool,
) {
	index.tree.Search(min, max, iter)
}

// Delete an item from the index
func (index *Index) Delete(min, max [2]float64, data interface{}) {
	index.tree.Delete(min, max, data)
}

// Children returns all children for parent node. If parent node is nil
// then the root nodes should be returned.
// The reuse buffer is an empty length slice that can optionally be used
// to avoid extra allocations.
func (index *Index) Children(parent interface{}, reuse []child.Child) (
	children []child.Child,
) {
	return index.tree.Children(parent, reuse)
}

// Nearby performs a kNN-type operation on the index.
// It's expected that the caller provides its own the `algo` function, which
// is used to calculate a distance to data. The `add` function should be
// called by the caller to "return" the data item along with a distance.
// The `iter` function will return all items from the smallest dist to the
// largest dist.
// Take a look at the SimpleBoxAlgo function for a usage example.
func (index *Index) Nearby(
	algo func(min, max [2]float64, data interface{}, item bool) (dist float64),
	iter func(min, max [2]float64, data interface{}, dist float64) bool,
) {
	var q queue
	var parent interface{}
	var children []child.Child

	for {
		// gather all children for parent
		children = index.tree.Children(parent, children[:0])
		for _, child := range children {
			q.push(qnode{
				dist:  algo(child.Min, child.Max, child.Data, child.Item),
				child: child,
			})
		}
		for {
			node, ok := q.pop()
			if !ok {
				// nothing left in queue
				return
			}
			if node.child.Item {
				if !iter(node.child.Min, node.child.Max,
					node.child.Data, node.dist) {
					return
				}
			} else {
				// gather more children
				parent = node.child.Data
				break
			}
		}
	}
}

// Len returns the number of items in tree
func (index *Index) Len() int {
	return index.tree.Len()
}

// Bounds returns the minimum bounding box
func (index *Index) Bounds() (min, max [2]float64) {
	return index.tree.Bounds()
}

// Scan iterates through all data in tree in no specified order.
func (index *Index) Scan(
	iter func(min, max [2]float64, data interface{}) bool,
) {
	index.tree.Scan(iter)
}

func (index *Index) svg(child child.Child, height int) []byte {
	var out []byte

	if !child.Item {
		out = append(out, fmt.Sprintf(
			"<rect x=\"%.0f\" y=\"%.0f\" width=\"%.0f\" height=\"%.0f\" "+
				"stroke=\"%s\" fill-opacity=\"0\" stroke-opacity=\"1\"/>\n",
			(child.Min[0])*svgScale,
			(child.Min[1])*svgScale,
			(child.Max[0]-child.Min[0]+1/svgScale)*svgScale,
			(child.Max[1]-child.Min[1]+1/svgScale)*svgScale,
			strokes[height%len(strokes)])...)
		children := index.tree.Children(child.Data, nil)
		for _, child := range children {
			out = append(out, index.svg(child, height+1)...)
		}
	} else {
		out = append(out, fmt.Sprintf(
			"<rect x=\"%.0f\" y=\"%.0f\" width=\"%.0f\" height=\"%.0f\" "+
				"stroke=\"%s\" fill-opacity=\"0\" stroke-opacity=\"1\"/>\n",
			(child.Min[0])*svgScale,
			(child.Min[1])*svgScale,
			(child.Max[0]-child.Min[0]+1/svgScale)*svgScale,
			(child.Max[1]-child.Min[1]+1/svgScale)*svgScale,
			strokes[len(strokes)-1])...)
	}
	return out
}

const svgScale = 5.0

var strokes = [...]string{"purple", "red", "#009900", "#cccc00", "black"}

// SVG prints 2D rtree in wgs84 coordinate space
func (index *Index) SVG() string {
	var out string
	out += fmt.Sprintf("<svg viewBox=\"%.0f %.0f %.0f %.0f\" "+
		"xmlns =\"http://www.w3.org/2000/svg\">\n",
		-190.0*svgScale, -100.0*svgScale,
		380.0*svgScale, 190.0*svgScale)

	out += fmt.Sprintf("<g transform=\"scale(1,-1)\">\n")

	var outb []byte
	for _, child := range index.Children(nil, nil) {
		outb = append(outb, index.svg(child, 1)...)
	}

	out += string(outb)
	out += fmt.Sprintf("</g>\n")
	out += fmt.Sprintf("</svg>\n")
	return out
}
