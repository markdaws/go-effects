package effects

import (
	"fmt"
	"image"
	"math"
)

// Rect used for image bounds
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// String returns a debug string
func (r Rect) String() string {
	return fmt.Sprintf("X:%d, Y:%d, Width:%d, Height:%d", r.X, r.Y, r.Width, r.Height)
}

// Intersect returns the intersection between two rectangles
func (r Rect) Intersect(r2 Rect) Rect {
	x := math.Max(float64(r.X), float64(r2.X))
	num1 := math.Min(float64(r.X+r.Width), float64(r2.X+r2.Width))

	y := math.Max(float64(r.Y), float64(r2.Y))
	num2 := math.Min(float64(r.Y+r.Height), float64(r2.Y+r2.Height))
	if num1 >= x && num2 >= y {
		return Rect{X: int(x), Y: int(y), Width: int(num1 - x), Height: int(num2 - y)}
	}
	return Rect{}
}

// IsEmpty returns true if this is an empty rectangle
func (r Rect) IsEmpty() bool {
	return r.Width == 0 || r.Height == 0
}

// ToImageRect returns an image.Rectangle instance initialized from this Rect
func (r Rect) ToImageRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: r.X, Y: r.Y},
		Max: image.Point{X: r.X + r.Width, Y: r.Y + r.Height},
	}
}
