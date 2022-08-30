package boids

import (
	"fmt"
	"math"
)

// Vector represents a 2D vector.
type Vector struct {
	X, Y float64
}

func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

// Calculates the vector angle and returns radians.
// To get degrees: multiply radians with 180/Pi
func (v Vector) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Calculates the dot product of two vectors.
func (v Vector) Dot(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Calculates the vector length/magnitude.
// It's an expensive call so use with care!
func (v Vector) Length() float64 {
	// Using math.Pow instead of plain x*x ensures consistent
	// and better rounding behaviours, but is waaay slower!
	// return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
	return math.Sqrt(v.Dot(v))
}

// Checks if vector length is within a target range.
// WARNING: target range r should be squared (r^2) by the caller!
// This odd behaviour is required for cutting down on the amount of square/square-roots needed within this function
// and hence provides some nice performance boost.
func (v Vector) InRange(r float64) float64 {
	d := v.Dot(v)
	if d < r {
		return math.Sqrt(d)
	}
	return 0
}

// Checks if the current vector is within a bounding box.
func (v Vector) Within(min, max Vector) bool {
	return v.X >= min.X && v.Y >= min.Y && v.X <= max.X && v.Y <= max.Y
}

func (v Vector) String() string {
	return fmt.Sprintf("(%+0.3f, %+0.3f)", v.X, v.Y)
}

// 6 decimal digits
const floatPrecision float64 = 1000000

func roundFloat(f float64) float64 {
	return math.Round(f*floatPrecision) / floatPrecision
}

func (v Vector) Round() Vector {
	v.X, v.Y = roundFloat(v.X), roundFloat(v.Y)
	return v
}

func (v Vector) Add(f float64) Vector {
	v.X += f
	v.Y += f
	return v
}

func (v Vector) Sub(f float64) Vector {
	v.X -= f
	v.Y -= f
	return v
}

func (v Vector) Mul(f float64) Vector {
	v.X *= f
	v.Y *= f
	return v
}

func (v Vector) Div(f float64) Vector {
	v.X /= f
	v.Y /= f
	return v
}

func (v Vector) Addv(other Vector) Vector {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v Vector) Subv(other Vector) Vector {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v Vector) Mulv(other Vector) Vector {
	v.X *= other.X
	v.Y *= other.Y
	return v
}

func (v Vector) Divv(other Vector) Vector {
	v.X /= other.X
	v.Y /= other.Y
	return v
}
