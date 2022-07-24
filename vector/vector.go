package vector

import (
	"fmt"
	"math"
)

const precision int = 6

type V struct {
	X, Y float64
}

func New(x, y float64) V {
	return V{x, y}
}

func (v V) String() string {
	return fmt.Sprintf("(%v, %v)", v.X, v.Y)
}

func (v V) Angle() float64 {
	// Calculates the vector angle and returns radians.
	// To get degrees: multiply radians with 180/Pi
	return math.Atan2(v.Y, v.X)
}

func (v V) Length() float64 {
	// Using math.Pow instead of plain x*x ensures consistent
	// and better rounding behaviours, but is waaay slower!
	//return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v V) Distance(other V) float64 {
	x, y := (v.X - other.X), (v.Y - other.Y)
	return math.Sqrt(x*x + y*y)
}

//func (v V) Equal(x, y float64) bool {
//return v.X == x && v.Y == x
//}

//func (v V) Equalv(other V) bool {
//return v.X == other.X && v.Y == other.Y
//}

//func (v V) Less(x, y float64) bool {
//return v.X < x && v.Y < y
//}

func (v V) Zero() V {
	v.X, v.Y = 0, 0
	return v
}

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
