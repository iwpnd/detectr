rtreego
=======

[![GoDoc](https://godoc.org/github.com/patrick-higgins/rtreego?status.svg)](https://godoc.org/github.com/patrick-higgins/rtreego)

A library for efficiently storing and querying spatial data
in the Go programming language.

Forked from github.com/dhconnelly/rtreego to specialize for ~~3~~ 2 dimensions
and tune for fewer memory allocations.


About
-----

The R-tree is a popular data structure for efficiently storing and
querying spatial objects; one common use is implementing geospatial
indexes in database management systems.  The variant implemented here,
known as the R*-tree, improves performance and increases storage
utilization.  Both bounding-box queries and k-nearest-neighbor queries
are supported.

R-trees are balanced, so maximum tree height is guaranteed to be
logarithmic in the number of entries; however, good worst-case
performance is not guaranteed.  Instead, a number of rebalancing
heuristics are applied that perform well in practice.  For more
details please refer to the references.


Status
------

Geometric primitives (points, rectangles, and their relevant geometric
algorithms) are implemented and tested.  The R-tree data structure and
algorithms are currently under development.


Install
-------

With Go 1 installed, just run `go get github.com/patrick-higgins/rtreego`.


Usage
-----

Make sure you `import github.com/patrick-higgins/rtreego` in your Go source files.

### Storing, updating, and deleting objects

To create a new tree, specify the number of spatial dimensions and the minimum
and maximum branching factor:

	rt := rtreego.NewTree(2, 25, 50)

Any type that implements the `Spatial` interface can be stored in the tree:

	type Spatial interface {
		Bounds() *Rect
	}

`Rect`s are data structures for representing spatial objects, while `Point`s
represent spatial locations.  Creating `Point`s is easy--they're just slices
of `float64`s:

	p1 := rtreego.Point{0.4, 0.5}
	p2 := rtreego.Point{6.2, -3.4}

To create a `Rect`, specify a location and the lengths of the sides:

	r1 := rtreego.NewRect(p1, []float64{1, 2})
	r2 := rtreego.NewRect(p2, []float64{1.7, 2.7})

To demonstrate, let's create and store some test data.

	type Thing struct {
		where *Rect
		name string
	}
	
	func (t *Thing) Bounds() *Rect {
		return t.where
	}
	
	rt.Insert(&Thing{r1, "foo"})
	rt.Insert(&Thing{r2, "bar"})
	
	size := rt.Size() // returns 2

We can insert and delete objects from the tree in any order.

	rt.Delete(thing2)
	// do some stuff...
	rt.Insert(anotherThing)

If you want to store points instead of rectangles, you can easily convert a
point into a rectangle using the `ToRect` method:

	var tol = 0.01

	type Somewhere struct {
		location rtreego.Point
		name string
		wormhole chan int
	}
	
	func (s *Somewhere) Bounds() *Rect {
		// define the bounds of s to be a rectangle centered at s.location
		// with side lengths 2 * tol:
		return s.location.ToRect(tol)
	}
	
	rt.Insert(&Somewhere{rtreego.Point{0, 0}, "Someplace", nil})

If you want to update the location of an object, you must delete it, update it,
and re-insert.  Just modifying the object so that the `*Rect` returned by 
`Location()` changes, without deleting and re-inserting the object, will
corrupt the tree.

### Queries

Bounding-box and k-nearest-neighbors queries are supported.

Bounding-box queries require a search `*Rect` argument and come in two flavors:
containment search and intersection search.  The former returns all objects that
fall strictly inside the search rectangle, while the latter returns all objects
that touch the search rectangle.

	bb := rtreego.NewRect(rtreego.Point{1.7, -3.4}, []float64{3.2, 1.9})

	// Get a slice of the objects in rt that intersect bb:
	results, _ := rt.SearchIntersect(bb)

	// Get a slice of the objects in rt that are contained inside bb:
	results, _ = rt.SearchContained(bb)

Nearest-neighbor queries find the objects in a tree closest to a specified
query point.

	q := rtreego.Point{6.5, -2.47}
	k := 5

	// Get a slice of the k objects in rt closest to q:
	results, _ = rt.SearchNearestNeighbors(q, k)

### More information

See http://github.com/patrick-higgins/rtreego for full API documentation.


References
----------

- A. Guttman.  R-trees: A Dynamic Index Structure for Spatial Searching.
  Proceedings of ACM SIGMOD, pages 47-57, 1984.
  http://www.cs.jhu.edu/~misha/ReadingSeminar/Papers/Guttman84.pdf
  
- N. Beckmann, H .P. Kriegel, R. Schneider and B. Seeger.  The R*-tree: An
  Efficient and Robust Access Method for Points and Rectangles.  Proceedings
  of ACM SIGMOD, pages 323-331, May 1990.
  http://infolab.usc.edu/csci587/Fall2011/papers/p322-beckmann.pdf
  
- N. Roussopoulos, S. Kelley and F. Vincent.  Nearest Neighbor Queries.  ACM
  SIGMOD, pages 71-79, 1995.
  http://www.postgis.org/support/nearestneighbor.pdf


Author
------

rtreego is written and maintained by Daniel Connelly.  You can find my stuff
at dhconnelly.com or email me at dhconnelly@gmail.com.

This fork is maintained by Patrick Higgins (patrick.allen.higgins@gmail.com).


License
-------
  
rtreego is released under a BSD-style license; see LICENSE for more details.
