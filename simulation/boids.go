package simulation

import (
	"math/rand"
	"sync"

	"github.com/lmas/boids/vector"
)

type Boid struct {
	Pos vector.V
	Vel vector.V
}

type Flock struct {
	Boids []*Boid
	Conf  Conf

	screenWidth  float64
	screenHeight float64
	center       vector.V
	wg           sync.WaitGroup
	signal       chan vector.V
	update       bool
}

func NewFlock(conf Conf) *Flock {
	f := &Flock{
		Boids:        make([]*Boid, conf.FlockSize),
		Conf:         conf,
		screenWidth:  float64(conf.ScreenWidth),
		screenHeight: float64(conf.ScreenHeight),
		signal:       make(chan vector.V, conf.GoRoutines),
	}
	f.center = vector.New(f.screenWidth/2, f.screenHeight/2)

	rand.Seed(conf.Seed)
	for i := 0; i < conf.FlockSize; i++ {
		f.Boids[i] = &Boid{
			Pos: vector.New(
				rand.Float64()*f.screenWidth,
				rand.Float64()*f.screenHeight,
			),
			Vel: vector.New(0, 0),
		}
	}

	// TODO: grab any leftovers if the flock wasn't divided up evenly
	group := f.Conf.FlockSize / f.Conf.GoRoutines
	for i := 0; i < f.Conf.GoRoutines; i++ {
		boids := f.Boids[i*group : (i*group)+group]
		go f.group(boids)
	}
	return f
}

func (f *Flock) Init(simulationSteps int) {
	for i := 0; i < simulationSteps; i++ {
		f.Step(true)
	}
}

func (f *Flock) Step(update bool) {
	f.update = update
	f.wg.Add(f.Conf.GoRoutines)
	for i := 0; i < f.Conf.GoRoutines; i++ {
		f.signal <- f.center
	}
	f.wg.Wait()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (f *Flock) group(boids []*Boid) {
	for {
		// TODO: check for termination signal so it can shut down cleanly?
		target := <-f.signal
		for _, b := range boids {
			f.stepBoid(b, target)
		}
		f.wg.Done()
	}
}

func (f *Flock) stepBoid(b *Boid, target vector.V) {
	sum := 0.0
	vel := vector.New(0, 0)
	cohesion := vector.New(0, 0)
	alignment := vector.New(0, 0)
	separation := vector.New(0, 0)

	// TODO: there's a race condition in here, should protect reading neighbour data with a lock.
	// But it seems to run just fine without, so far?
	if f.update {
		for _, n := range f.Boids {
			if n == b || b.Pos.Distance(n.Pos) > f.Conf.VisionRadious {
				continue
			}
			sum += 1.0
			cohesion = cohesion.Addv(n.Pos)
			alignment = alignment.Addv(n.Vel)
			if b.Pos.Distance(n.Pos) < f.Conf.SeparationRadious {
				separation = separation.Addv(b.Pos.Subv(n.Pos))
				if separation.Length() < 1 {
					separation = vector.New(rand.Float64(), rand.Float64())
				}
			}
		}
		if sum > 0 {
			cohesion = cohesion.Div(sum).Subv(b.Pos)
			alignment = alignment.Div(sum).Subv(b.Vel)
			vel = vel.Addv(cohesion.Mul(f.Conf.CohesionFactor))
			vel = vel.Addv(alignment.Mul(f.Conf.AlignmentFactor))
			vel = vel.Addv(separation.Mul(f.Conf.SeparationFactor))
		}
	}

	vel = vel.Addv(target.Subv(b.Pos).Mul(f.Conf.TargetingFactor))
	b.Vel = f.limitSpeed(b.Vel.Addv(vel)).Round()
	b.Pos = b.Pos.Addv(b.Vel).Round()
}

func (f *Flock) limitSpeed(v vector.V) vector.V {
	len := v.Length()
	if len > f.Conf.SpeedLimitingFactor {
		v = v.Div(len)
	}
	return v.Mul(f.Conf.SpeedLimitingFactor)
}
