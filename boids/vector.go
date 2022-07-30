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

func (v Vector) String() string {
	return fmt.Sprintf("(%+0.3f, %+0.3f)", v.X, v.Y)
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

// A.K.A. Magnitude
func (v Vector) Length() float64 {
	// Using math.Pow instead of plain x*x ensures consistent
	// and better rounding behaviours, but is waaay slower!
	// return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
	return math.Sqrt(v.Dot(v))
}

// Checks if the current vector is within a bounding box.
func (v Vector) Within(min, max Vector) bool {
	return v.X >= min.X && v.Y >= min.Y && v.X <= max.X && v.Y <= max.Y
}

const vectorPrecision int = 6

func (v Vector) Round() Vector {
	output := math.Pow(10, float64(vectorPrecision))
	v.X = math.Round(v.X*output) / output
	v.Y = math.Round(v.Y*output) / output
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
