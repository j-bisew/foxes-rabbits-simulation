package quadtree

import (
	"fmt"
	"strings"

	"github.com/j-bisew/foxes-rabbits-simulation/geom"
)


type Point = geom.Point
type Rectangle = geom.Rectangle

type QuadTree struct {
	capacity int
	boundary Rectangle
	points []Point
	nw,ne,sw,se *QuadTree
	divided bool
}

func NewQuadTree(capacity int, rect Rectangle) *QuadTree {
	return &QuadTree{
		capacity: capacity,
		boundary: rect,
		points:   make([]Point, 0, capacity),
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

	for _, point := range qt.points {
		_ = qt.ne.Insert(point) ||
		qt.nw.Insert(point) ||
		qt.se.Insert(point) ||
		qt.sw.Insert(point)
	}

	qt.points = nil
}

func (qt *QuadTree) Insert(point Point) bool {
	if !qt.boundary.Contains(point) {
		return false
	}

	if !qt.divided {
		if len(qt.points) < qt.capacity {
			qt.points = append(qt.points, point)
			return true
		}
		qt.subdivide()
	}
	return qt.ne.Insert(point) ||
		qt.nw.Insert(point) ||
		qt.se.Insert(point) ||
		qt.sw.Insert(point)
}

func (qt *QuadTree) Query(rangeRect Rectangle, found *[]Point) {
	if !qt.boundary.Intersects(rangeRect) {
		return
	}
	
	if !qt.divided {
		for _, p := range qt.points {
			if rangeRect.Contains(p) {
				*found = append(*found, p)
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
				sb.WriteString(fmt.Sprintf("%s [%d] (%.1f, %.1f)\n", indent, i, point.X, point.Y))
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
