package boids

import (
	"math/rand"
	"sync"
)

type Conf struct {
	Seed        int64
	Workers     int
	SwarmSize   int
	IndexOffset int
}

// Swarm is a group of Boids.
// It is moving together most of the time.
type Swarm struct {
	Conf  Conf
	Boids []*Boid
	Index *Index

	signal chan workerSignal
	wg     sync.WaitGroup
}

// New creates a new swarm of Boids, using Conf.
// It randomises the positions of each Boid and fires up a group of background
// workers to perform the actual Boid movement updates.
func New(conf Conf) *Swarm {
	s := &Swarm{
		Conf:   conf,
		Boids:  make([]*Boid, conf.SwarmSize),
		Index:  NewIndex(conf.IndexOffset),
		signal: make(chan workerSignal, conf.Workers),
	}

	rand.Seed(conf.Seed)
	for i := 0; i < conf.SwarmSize; i++ {
		s.Boids[i] = &Boid{
			ID:  i,
			Pos: NewVector(rand.Float64()-0.5, rand.Float64()-0.5).Mul(10),
			Vel: NewVector(0, 0),
		}
	}

	// TODO: grab any leftovers if the flock wasn't divided up evenly
	worker := conf.SwarmSize / conf.Workers
	for i := 0; i < conf.Workers; i++ {
		boids := s.Boids[i*worker : (i*worker)+worker]
		go s.workerUpdate(boids)
	}
	return s
}

// Update all Boids' velocity (dirty, slow) or position (non-dirty, fast).
// It also updates the Boid neighbour index if dirty, before hand.
func (s *Swarm) Update(dirty bool, target Vector) {
	// TODO: could allow multiple targets?
	if dirty {
		s.Index.Update(s.Boids)
	}

	sig := workerSignal{dirty, target}
	s.wg.Add(s.Conf.Workers)
	for i := 0; i < s.Conf.Workers; i++ {
		s.signal <- sig
	}
	s.wg.Wait()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type workerSignal struct {
	Dirty  bool
	Target Vector
}

func (s *Swarm) workerUpdate(boids []*Boid) {
	for {
		// TODO: check for termination signal so it can shut down cleanly?
		sig := <-s.signal
		for _, b := range boids {
			s.updateBoid(b, sig.Dirty, sig.Target)
		}
		s.wg.Done()
	}
}
