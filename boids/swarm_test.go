package boids

import (
	"testing"

	"github.com/lmas/boids/vector"
)

func BenchmarkBoids(b *testing.B) {
	s := New(Conf{
		Seed:       0,
		GoRoutines: 10,
		SwarmSize:  500,
	})
	v := vector.New(0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Must alternate between updating velocity (dirty) and position (non-dirty)
		s.Update(i%2 == 0, v)
	}
}
