package boids

import (
	"fmt"
	"math"
	"testing"
)

const pi float64 = math.Pi

var piv = NewVector(math.Pi, math.Pi)

func assertVector(t *testing.T, v Vector, x, y float64) {
	if v.X != x || v.Y != y {
		t.Errorf("got vector %s, expected (%v, %v)", v, x, y)
	}
}

func TestSimpleOperations(t *testing.T) {
	t.Run("angle", func(t *testing.T) {
		f := piv.Angle()
		e := math.Atan2(pi, pi)
		if f != e {
			t.Errorf("got angle %f, expected %f", f, e)
		}
	})
	t.Run("dot product", func(t *testing.T) {
		f := piv.Dot(piv)
		e := pi*pi + pi*pi
		if f != e {
			t.Errorf("got dot product %f, expected %f", f, e)
		}
	})
	t.Run("vector length", func(t *testing.T) {
		f := piv.Length()
		e := math.Sqrt(pi*pi + pi*pi)
		if f != e {
			t.Errorf("got vector length %f, expected %f", f, e)
		}
	})
	t.Run("within", func(t *testing.T) {
		min := NewVector(0, 0)
		assertVector(t, min, 0, 0)
		max := NewVector(pi, pi)
		assertVector(t, max, pi, pi)
		b := piv.Within(min, max)
		e := true
		if b != e {
			t.Errorf("got vector within '%v', expected '%v'", b, e)
		}
	})
	t.Run("string", func(t *testing.T) {
		s := piv.String()
		e := fmt.Sprintf("(%+0.3f, %+0.3f)", pi, pi)
		if s != e {
			t.Errorf("got vector string '%s', expected '%s'", s, e)
		}
	})
	t.Run("round", func(t *testing.T) {
		assertVector(t, piv, pi, pi)
		v := piv.Round()
		e := 3.141593
		assertVector(t, v, e, e)
	})
	t.Run("float arithmetic", func(t *testing.T) {
		v := NewVector(0, 0)
		assertVector(t, v, 0, 0)
		v = v.Add(pi)
		assertVector(t, v, pi, pi)
		v = v.Mul(pi)
		assertVector(t, v, pi*pi, pi*pi)
		v = v.Div(pi)
		assertVector(t, v, pi, pi)
		v = v.Sub(pi)
		assertVector(t, v, 0, 0)
	})
	t.Run("vector arithmetic", func(t *testing.T) {
		v := NewVector(0, 0)
		assertVector(t, v, 0, 0)
		v = v.Addv(piv)
		assertVector(t, v, pi, pi)
		v = v.Mulv(piv)
		assertVector(t, v, pi*pi, pi*pi)
		v = v.Divv(piv)
		assertVector(t, v, pi, pi)
		v = v.Subv(piv)
		assertVector(t, v, 0, 0)
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkVectors(b *testing.B) {
	test := NewVector(0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		test = test.Add(pi)
		test = test.Mul(pi)
		test = test.Div(pi)
		test = test.Sub(pi)

		test = test.Addv(piv)
		test = test.Mulv(piv)
		test = test.Divv(piv)
		test = test.Subv(piv)

		test.Round()
	}
}
