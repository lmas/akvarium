package simulation

import (
	"testing"

	"github.com/lmas/boids/vector"
)

var conf = DefaultConf()

func BenchmarkBoids(b *testing.B) {
	s := NewSwarm(conf)
	v := vector.New(0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Update(true, v)
	}
}
