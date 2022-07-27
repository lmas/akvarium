package simulation

import "testing"

var conf = DefaultConf()

func BenchmarkBoids(b *testing.B) {
	s := NewSwarm(conf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Update(true)
	}
}
