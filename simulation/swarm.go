package simulation

import (
	"math/rand"
	"sync"

	"github.com/lmas/boids/vector"
)

type groupSignal struct {
	Dirty  bool
	Target vector.V
}

type Swarm struct {
	Conf   Conf
	Boids  []*Boid
	index  binIndex
	signal chan groupSignal
	wg     sync.WaitGroup
}

func NewSwarm(conf Conf) *Swarm {
	s := &Swarm{
		Conf:   conf,
		Boids:  make([]*Boid, conf.SwarmSize),
		signal: make(chan groupSignal, conf.GoRoutines),
	}

	rand.Seed(conf.Seed)
	for i := 0; i < conf.SwarmSize; i++ {
		s.Boids[i] = &Boid{
			ID:  i,
			Pos: vector.New(rand.Float64()-0.5, rand.Float64()-0.5).Mul(10),
			Vel: vector.New(0, 0),
		}
	}

	// TODO: grab any leftovers if the flock wasn't divided up evenly
	group := conf.SwarmSize / conf.GoRoutines
	for i := 0; i < conf.GoRoutines; i++ {
		boids := s.Boids[i*group : (i*group)+group]
		go s.updateGroup(boids)
	}
	return s
}

func (s *Swarm) Update(dirty bool, target vector.V) {
	// TODO: could allow multiple targets?
	if dirty {
		s.Index()
	}

	sig := groupSignal{dirty, target}
	s.wg.Add(s.Conf.GoRoutines)
	for i := 0; i < s.Conf.GoRoutines; i++ {
		s.signal <- sig
	}
	s.wg.Wait()
}

func (s *Swarm) updateGroup(boids []*Boid) {
	for {
		// TODO: check for termination signal so it can shut down cleanly?
		sig := <-s.signal
		for _, b := range boids {
			s.updateBoid(b, sig.Dirty, sig.Target)
		}
		s.wg.Done()
	}
}
