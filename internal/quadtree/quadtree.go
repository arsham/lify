// Package quadtree implements a quadtree datastructure for storing points in a
// 2D space for fast lookup.
package quadtree

import (
	"fmt"
	"strings"
	"sync"

	"github.com/arsham/neuragene/internal/geom"
)

type quadrant uint8

const (
	// NW is the North West quadrant.
	NW quadrant = iota
	// NE is the North East quadrant.
	NE
	// SW is the South West quadrant.
	SW
	// SE is the South East quadrant.
	SE
)

// A Point is a point in 2D space for keeping track of entities on the
// quadtree.
type Point[T any] struct {
	Data T
	geom.Vec
}

// P creates a new Point object from the given data and vector.
func P[T any](data T, v geom.Vec) Point[T] {
	return Point[T]{
		Data: data,
		Vec:  v,
	}
}

func (p Point[T]) String() string {
	return fmt.Sprintf("Point(x: %.1f, y: %.1f, data: %v)", p.X, p.Y, p.Data)
}

// Bounds is a rectangle that represents the boundaries of a quadtree node. The
// Min is the top left corner, and the Max is the bottom right corner.
type Bounds struct {
	geom.Rect
}

var boundsPool = sync.Pool{
	New: func() interface{} {
		return &Bounds{}
	},
}

// NewBounds creates a new Bounds object. The x1 and y1 are the coordinates of
// the top left corner of the bounds, and the x2 and y2 are the coordinates of
// the bottom right corner of the bounds.
func NewBounds(x1, y1, x2, y2 float64) *Bounds {
	// nolint:forcetypeassert // we already know the type.
	b := boundsPool.Get().(*Bounds)
	b.Min.X = x1
	b.Min.Y = y1
	b.Max.X = x2
	b.Max.Y = y2
	return b
}

// Free returns the Bounds to the pool.
func (b *Bounds) Free() {
	boundsPool.Put(b)
}

// Contains returns true if the point is within the bounds. If the point is
// exactly on the edge of the bounds, it is considered to be within the bounds.
func (b *Bounds) Contains(point geom.Vec) bool {
	return b.Rect.Contains(point)
}

// Intersects returns true if the bounds intersects with the other bounds. If
// either bounds is exactly on the edge of the other bounds, it is considered
// to be an intersection.
func (b *Bounds) Intersects(box geom.Rect) bool {
	return b.Rect.Intersects(box)
}

// SubDivide returns a new Bounds object that represents the quadrant of the
// current bounds. It panics of the q quadrant is not one of the four defined
// quadrants.
func (b *Bounds) SubDivide(q quadrant) *Bounds {
	switch q {
	case NW:
		return NewBounds(
			b.Min.X,
			b.Min.Y,
			b.Min.X+b.W()/2,
			b.Min.Y+b.H()/2,
		)
	case NE:
		return NewBounds(
			b.Min.X+b.W()/2,
			b.Min.Y,
			b.Max.X,
			b.Min.Y+b.H()/2,
		)
	case SW:
		return NewBounds(
			b.Min.X,
			b.Min.Y+b.H()/2,
			b.Min.X+b.W()/2,
			b.Max.Y,
		)
	case SE:
		return NewBounds(
			b.Min.X+b.W()/2,
			b.Min.Y+b.H()/2,
			b.Max.X,
			b.Max.Y,
		)
	}
	panic(fmt.Sprintf("invalid quadrant: %d", q))
}

var quadTreePool = sync.Pool{
	New: func() interface{} {
		return &QuadTree[uint64]{
			points: make([]Point[uint64], 0, 1),
		}
	},
}

// QuadTree holds a series of points in a quadtree structure. If might contain
// four children, and each child might contain four children, and so on. The
// top level quadtree is the root node. When the number of points in a node
// exceeds the capacity, the node is subdivided into four children. The points
// are then moved to the smallest possible rectangle. This improves performance
// by reducing the number of points that need to be checked.
type QuadTree[T any] struct {
	boundary *Bounds
	nw       *QuadTree[T]
	ne       *QuadTree[T]
	sw       *QuadTree[T]
	se       *QuadTree[T]
	points   []Point[T]
	capacity uint
	depth    uint8 // used for debugging purposes.
	divided  bool
}

// NewQuadTree creates a new QuadTree object with the given boundary and
// maximum capacity for each quadrant.
func NewQuadTree[T any](b *Bounds, capacity uint, depth uint8) *QuadTree[T] {
	// nolint:forcetypeassert // we already know the type.
	q := quadTreePool.Get().(*QuadTree[T])
	q.boundary = b
	q.capacity = capacity
	q.depth = depth
	return q
}

// Free returns the QuadTree to the pool.
func (q *QuadTree[T]) Free() {
	if q.nw != nil {
		q.nw.Free()
		q.nw = nil
	}
	if q.ne != nil {
		q.ne.Free()
		q.ne = nil
	}
	if q.sw != nil {
		q.sw.Free()
		q.sw = nil
	}
	if q.se != nil {
		q.se.Free()
		q.se = nil
	}
	q.points = q.points[:0]
	q.depth = 0
	q.divided = false
	q.boundary.Free()
	quadTreePool.Put(q)
}

// Children returns the children of the quadtree node. It returns nil if this
// node is not divided yet.
func (q *QuadTree[T]) Children() []*QuadTree[T] {
	if !q.divided {
		return nil
	}
	return []*QuadTree[T]{
		q.nw,
		q.ne,
		q.sw,
		q.se,
	}
}

// SubDivide creates four new quadtree nodes and moves the points to the new
// nodes.
func (q *QuadTree[T]) SubDivide() {
	q.nw = NewQuadTree[T](
		q.boundary.SubDivide(NW),
		q.capacity,
		q.depth+1,
	)
	q.ne = NewQuadTree[T](
		q.boundary.SubDivide(NE),
		q.capacity,
		q.depth+1,
	)
	q.sw = NewQuadTree[T](
		q.boundary.SubDivide(SW),
		q.capacity,
		q.depth+1,
	)
	q.se = NewQuadTree[T](
		q.boundary.SubDivide(SE),
		q.capacity,
		q.depth+1,
	)
	q.divided = true

	// Move points to children. This improves performance by placing points in
	// the smallest available rectangle.
	for _, p := range q.points {
		q.nw.Insert(p)
		q.ne.Insert(p)
		q.sw.Insert(p)
		q.se.Insert(p)
	}
	q.points = nil
}

// Insert adds a point to the quadtree. If the capacity is exceeded, the node
// is subdivided into four children. If the point is not within the bounds of
// this node, it is not added.
func (q *QuadTree[T]) Insert(p Point[T]) {
	if !q.boundary.Contains(p.Vec) {
		return
	}
	if !q.divided {
		if len(q.points) < int(q.capacity) {
			q.points = append(q.points, p)
			return
		}
		q.SubDivide()
	}
	q.nw.Insert(p)
	q.ne.Insert(p)
	q.sw.Insert(p)
	q.se.Insert(p)
}

var queryResultPool = sync.Pool{
	New: func() interface{} {
		return &QueryResult[uint64]{
			Points: make([]Point[uint64], 0, 1000),
		}
	},
}

// QueryResult is a result of a query.
type QueryResult[T any] struct {
	Points []Point[T]
}

// Free returns the QueryResult to the pool.
func (q *QueryResult[T]) Free() {
	q.Points = q.Points[:0]
	queryResultPool.Put(q)
}

// Q returns a new QueryResult object from a pool of QueryResult objects.
func Q[T any]() *QueryResult[T] {
	// nolint:forcetypeassert // we already know the type.
	return queryResultPool.Get().(*QueryResult[T])
}

// Query returns all points within the bounds.
func (q *QuadTree[T]) Query(rect geom.Rect, result *QueryResult[T]) {
	if !q.boundary.Intersects(rect) {
		return
	}
	if q.divided {
		q.nw.Query(rect, result)
		q.ne.Query(rect, result)
		q.sw.Query(rect, result)
		q.se.Query(rect, result)
		return
	}

	for _, p := range q.points {
		if rect.Contains(p.Vec) {
			result.Points = append(result.Points, p)
		}
	}
}

// Bounds returns the top left x and y, and bottom right x and y.
func (q *QuadTree[T]) Bounds() (x1, y1, x2, y2 float64) {
	return q.boundary.Min.X, q.boundary.Min.Y, q.boundary.Max.X, q.boundary.Max.Y
}

func (q *QuadTree[T]) String() string {
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "%sQT(Min.X: %.2f, Min.Y: %.2f, Max.X: %.2f, Max.Y: %.2f)",
		strings.Repeat("  ", int(q.depth)), q.boundary.Min.X, q.boundary.Min.Y, q.boundary.Max.X, q.boundary.Max.Y)
	for _, c := range q.Children() {
		fmt.Fprintf(builder, "\n%s", c)
	}
	return builder.String()
}
