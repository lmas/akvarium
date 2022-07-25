package simulation

import "testing"

var conf = DefaultConf()

func BenchmarkBoids(b *testing.B) {
	f := NewFlock(conf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Step(true)
	}
}
