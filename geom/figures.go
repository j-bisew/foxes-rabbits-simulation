package geom

type Point struct {
	X, Y float64
}

type Rectangle struct {
	X, Width, Y, Height float64
}

func NewRectangle(x, width, y, height float64) *Rectangle {
	return &Rectangle{x, width, y, height}
}

func (r *Rectangle) Contains(p Point) bool {
	return p.X >= r.X && p.X <= r.X+r.Width && p.Y >= r.Y && p.Y <= r.Y+r.Height
}

func (r *Rectangle) Intersects(other Rectangle) bool {
	return !(r.X > other.X+other.Width || r.X+r.Width < other.X ||
	r.Y > other.Y+other.Height || r.Y+r.Height < other.Y)
}
