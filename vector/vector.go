package vector

import (
	"fmt"
	"math"
)

// V represents a 2D vector.
type V struct {
	X, Y float64
}

func New(x, y float64) V {
	return V{x, y}
}

func (v V) String() string {
	return fmt.Sprintf("(%+0.3f, %+0.3f)", v.X, v.Y)
}

// Calculates the vector angle and returns radians.
// To get degrees: multiply radians with 180/Pi
func (v V) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Calculates the dot product of two vectors.
func (v V) Dot(other V) float64 {
	return v.X*other.X + v.Y*other.Y
}

// A.K.A. Magnitude
func (v V) Length() float64 {
	// Using math.Pow instead of plain x*x ensures consistent
	// and better rounding behaviours, but is waaay slower!
	// return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
	return math.Sqrt(v.Dot(v))
}

// Checks if the current vector is within a bounding box.
func (v V) Within(min, max V) bool {
	return v.X >= min.X && v.Y >= min.Y && v.X <= max.X && v.Y <= max.Y
}

const precision int = 6

func (v V) Round() V {
	output := math.Pow(10, float64(precision))
	v.X = math.Round(v.X*output) / output
	v.Y = math.Round(v.Y*output) / output
	return v
}

func (v V) Add(f float64) V {
	v.X += f
	v.Y += f
	return v
}

func (v V) Sub(f float64) V {
	v.X -= f
	v.Y -= f
	return v
}

func (v V) Mul(f float64) V {
	v.X *= f
	v.Y *= f
	return v
}

func (v V) Div(f float64) V {
	v.X /= f
	v.Y /= f
	return v
}

func (v V) Addv(other V) V {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v V) Subv(other V) V {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v V) Mulv(other V) V {
	v.X *= other.X
	v.Y *= other.Y
	return v
}

func (v V) Divv(other V) V {
	v.X /= other.X
	v.Y /= other.Y
	return v
}
