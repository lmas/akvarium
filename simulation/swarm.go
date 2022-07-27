package simulation

import (
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lmas/boids/vector"
)

type Swarm struct {
	Conf   Conf
	Boids  []*Boid
	screen vector.V
	center vector.V
	target vector.V
	index  binIndex
	signal chan bool
	wg     sync.WaitGroup
}

func NewSwarm(conf Conf) *Swarm {
	s := &Swarm{
		Conf:   conf,
		Boids:  make([]*Boid, conf.SwarmSize),
		screen: vector.New(float64(conf.ScreenWidth), float64(conf.ScreenHeight)),
		center: vector.New(float64(conf.ScreenWidth)/2, float64(conf.ScreenHeight)/2),
		signal: make(chan bool, conf.GoRoutines),
	}

	rand.Seed(conf.Seed)
	for i := 0; i < conf.SwarmSize; i++ {
		s.Boids[i] = &Boid{
			ID:  i,
			Pos: vector.New(rand.Float64(), rand.Float64()).Mulv(s.screen),
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

func (s *Swarm) Update(dirty bool) {
	if dirty {
		s.Index()
		s.target = s.center
		cx, cy := ebiten.CursorPosition()
		cursor := vector.New(float64(cx), float64(cy))
		if cursor.Within(vector.New(0, 0), s.screen) {
			s.target = cursor
		}
	}

	s.wg.Add(s.Conf.GoRoutines)
	for i := 0; i < s.Conf.GoRoutines; i++ {
		s.signal <- dirty
	}
	s.wg.Wait()
}

func (s *Swarm) updateGroup(boids []*Boid) {
	for {
		// TODO: check for termination signal so it can shut down cleanly?
		dirty := <-s.signal
		for _, b := range boids {
			s.updateBoid(b, dirty, s.target)
		}
		s.wg.Done()
	}
}
