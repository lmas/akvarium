package boids

import (
	"math/rand"
	"sync"
)

type Conf struct {
	Boids       int
	Spawn       [2]Vector
	Seed        int64
	Workers     int
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
		Boids:  make([]*Boid, conf.Boids),
		Index:  NewIndex(conf.IndexOffset),
		signal: make(chan workerSignal, conf.Workers),
	}

	min, max := conf.Spawn[0], conf.Spawn[1]
	rand.Seed(conf.Seed)
	for i := 0; i < conf.Boids; i++ {
		s.Boids[i] = &Boid{
			ID: i,
			Pos: NewVector(
				min.X+rand.Float64()*(max.X-min.X), //nolint:gosec
				min.Y+rand.Float64()*(max.Y-min.Y), //nolint:gosec
			),
			Vel: NewVector(0, 0),
		}
	}

	// TODO: grab any leftovers if the flock wasn't divided up evenly
	worker := conf.Boids / conf.Workers
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
