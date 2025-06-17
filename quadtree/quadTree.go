package quadtree

import (
	"fmt"
	"strings"

	"github.com/j-bisew/foxes-rabbits-simulation/geom"
	"github.com/j-bisew/foxes-rabbits-simulation/interfaces"
)


type Point = geom.Point
type Rectangle = geom.Rectangle
type Entity = interfaces.Entity

type EntityPoint struct {
	Point Point
	Entity Entity
}

type QuadTree struct {
	capacity int
	boundary Rectangle
	points []EntityPoint
	nw,ne,sw,se *QuadTree
	divided bool
}

func NewQuadTree(capacity int, rect Rectangle) *QuadTree {
	return &QuadTree{
		capacity: capacity,
		boundary: rect,
		points:   make([]EntityPoint, 0, capacity),
		divided:  false,
	}
}

func (qt *QuadTree) subdivide() {
	x := qt.boundary.X
	y := qt.boundary.Y
	w := qt.boundary.Width / 2
	h := qt.boundary.Height / 2

	qt.ne = NewQuadTree(qt.capacity, geom.Rectangle{X: x+w, Y: y+h, Width: w, Height: h})
	qt.nw = NewQuadTree(qt.capacity, geom.Rectangle{X: x, Y: y+h, Width: w, Height: h})
	qt.se = NewQuadTree(qt.capacity, geom.Rectangle{X: x+w, Y: y, Width: w, Height: h})
	qt.sw = NewQuadTree(qt.capacity, geom.Rectangle{X: x, Y: y, Width: w, Height: h})

	qt.divided = true

	for _, ep := range qt.points {
		_ = qt.ne.Insert(ep.Point, ep.Entity) ||
		qt.nw.Insert(ep.Point, ep.Entity) ||
		qt.se.Insert(ep.Point, ep.Entity) ||
		qt.sw.Insert(ep.Point, ep.Entity)
	}

	qt.points = nil
}

func (qt *QuadTree) Clear() {
	qt.points = qt.points[:0]
	qt.divided = false
	qt.ne, qt.nw, qt.se, qt.sw = nil, nil, nil, nil
}

func (qt *QuadTree) Insert(point Point, entity Entity) bool {
	if !qt.boundary.Contains(point) {
		return false
	}

	entityPoint := EntityPoint{Point: point, Entity: entity}

	if !qt.divided {
		if len(qt.points) < qt.capacity {
			qt.points = append(qt.points, entityPoint)
			return true
		}
		qt.subdivide()
	}
	return qt.ne.Insert(point, entity) ||
		qt.nw.Insert(point, entity) ||
		qt.se.Insert(point, entity) ||
		qt.sw.Insert(point, entity)
}

func (qt *QuadTree) Query(rangeRect Rectangle, found *[]Entity) {
	if !qt.boundary.Intersects(rangeRect) {
		return
	}
	
	if !qt.divided {
		for _, ep := range qt.points {
			if rangeRect.Contains(ep.Point) {
				*found = append(*found, ep.Entity)
			}
		}
		return
	}
	
	qt.ne.Query(rangeRect, found)
	qt.nw.Query(rangeRect, found)
	qt.se.Query(rangeRect, found)
	qt.sw.Query(rangeRect, found)
}

func (qt *QuadTree) String() string {
	return qt.stringWithIndent("")
}

func (qt *QuadTree) stringWithIndent(indent string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%sQuadtree [X:%.1f, Y:%.1f  %.1fx%.1f] cap:%d", indent, qt.boundary.X, qt.boundary.Y, qt.boundary.Width, qt.boundary.Height, qt.capacity))

	if !qt.divided {
		if len(qt.points) > 0 {
			sb.WriteString(fmt.Sprintf(" points:%d\n", len(qt.points)))
			for i, point := range qt.points {
				sb.WriteString(fmt.Sprintf("%s [%d] (%.1f, %.1f)\n", indent, i, point.Point.X, point.Point.Y))
			}
		} else {
			sb.WriteString(" (empty)\n")
		}
	} else {
		sb.WriteString(" (subdivided)\n")
		childIndent := indent + "  "
		sb.WriteString(fmt.Sprintf("%sNE: %s", childIndent, qt.ne.stringWithIndent(childIndent)))
		sb.WriteString(fmt.Sprintf("%sNW: %s", childIndent, qt.nw.stringWithIndent(childIndent)))
		sb.WriteString(fmt.Sprintf("%sSE: %s", childIndent, qt.se.stringWithIndent(childIndent)))
		sb.WriteString(fmt.Sprintf("%sSW: %s", childIndent, qt.sw.stringWithIndent(childIndent)))
	}

	return sb.String()
}
