package boids

import (
	"math"
	"testing"
)

func assertFloat(t *testing.T, got, expected float64) {
	if got != expected {
		t.Errorf("got float %v, expected %v", got, expected)
	}
}

func assertVector(t *testing.T, v Vector, x, y float64) {
	if v.X != x || v.Y != y {
		t.Errorf("got vector %s, expected (%v, %v)", v, x, y)
	}
}

func TestSimpleOperations(t *testing.T) {
	v := NewVector(0, 0)
	t.Run("new is zero", func(t *testing.T) {
		assertVector(t, v, 0, 0)
	})
	t.Run("get angle", func(t *testing.T) {
		f := math.Atan2(3.3, 3.3)
		assertFloat(t, v.Add(3.3).Angle(), f)
	})
	t.Run("get length", func(t *testing.T) {
		f := math.Sqrt(math.Pow(3.3, 2) + math.Pow(3.3, 2))
		assertFloat(t, v.Add(3.3).Length(), f)
	})
	t.Run("do round", func(t *testing.T) {
		assertVector(t, v.Add(math.Sqrt(3.3*3.3)).Round(), 3.3, 3.3)
	})
	t.Run("add float", func(t *testing.T) {
		assertVector(t, v.Add(3.3), 3.3, 3.3)
	})
	t.Run("sub float", func(t *testing.T) {
		assertVector(t, v.Sub(3.3), -3.3, -3.3)
	})
	t.Run("mul float", func(t *testing.T) {
		assertVector(t, v.Add(1).Mul(3.3), 3.3, 3.3)
	})
	t.Run("div float", func(t *testing.T) {
		f := float64(1 / 3.3)
		assertVector(t, v.Add(1).Div(3.3), f, f)
	})
	t.Run("add vector", func(t *testing.T) {
		assertVector(t, v.Addv(NewVector(3.3, 3.3)), 3.3, 3.3)
	})
	t.Run("sub vector", func(t *testing.T) {
		assertVector(t, v.Subv(NewVector(3.3, 3.3)), -3.3, -3.3)
	})
	t.Run("mul vector", func(t *testing.T) {
		assertVector(t, v.Add(1).Mulv(NewVector(3.3, 3.3)), 3.3, 3.3)
	})
	t.Run("div vector", func(t *testing.T) {
		f := float64(1 / 3.3)
		assertVector(t, v.Add(1).Divv(NewVector(3.3, 3.3)), f, f)
	})
	t.Run("chain multiple operations", func(t *testing.T) {
		v = v.Add(4.4).Sub(1.1).Mul(2.2).Div(2.2).Round()
		assertVector(t, v, 3.3, 3.3)
	})
}
