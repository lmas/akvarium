package boids

import (
	"math"
)

type bin []int
type binIndex map[binKey]bin

type binKey [2]int

func getBinKey(b *Boid) binKey {
	v := b.Pos.Div(neighbourRange)
	return binKey{
		int(math.Floor(v.X)),
		int(math.Floor(v.Y)),
	}
}

func (s *Swarm) Index() {
	s.index = make(binIndex)
	for _, b := range s.Boids {
		k := getBinKey(b)
		s.index[k] = append(s.index[k], b.ID)
	}
}

func (s *Swarm) IterNeighbours(b *Boid, fun func(n *Boid)) {
	k := getBinKey(b)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			s.iterBin(binKey{k[0] + i, k[1] + j}, b, fun)
		}
	}
}

func (s *Swarm) iterBin(k binKey, b *Boid, fun func(n *Boid)) {
	for _, i := range s.index[k] {
		if i == b.ID {
			continue
		}
		fun(s.Boids[i])
	}
}
