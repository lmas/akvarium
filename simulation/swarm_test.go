package simulation

import (
	"testing"

	"github.com/lmas/boids/vector"
)

func BenchmarkBoids(b *testing.B) {
	s := NewSwarm(Conf{
		Seed:       0,
		GoRoutines: 10,
		SwarmSize:  500,
	})
	v := vector.New(0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Update(true, v)
	}
}
