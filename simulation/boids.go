package simulation

import (
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lmas/boids/vector"
)

type Boid struct {
	ID  int
	Pos vector.V
	Vel vector.V
}

type Flock struct {
	Boids []*Boid
	Conf  Conf

	screen vector.V
	center vector.V
	wg     sync.WaitGroup
	signal chan vector.V
	update bool
}

func NewFlock(conf Conf) *Flock {
	f := &Flock{
		Boids:  make([]*Boid, conf.FlockSize),
		Conf:   conf,
		screen: vector.New(float64(conf.ScreenWidth), float64(conf.ScreenHeight)),
		center: vector.New(float64(conf.ScreenWidth)/2, float64(conf.ScreenHeight)/2),
		signal: make(chan vector.V, conf.GoRoutines),
	}

	rand.Seed(conf.Seed)
	for i := 0; i < conf.FlockSize; i++ {
		f.Boids[i] = &Boid{
			ID:  i,
			Pos: vector.New(rand.Float64(), rand.Float64()).Mulv(f.screen),
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
	target := f.center
	if update {
		cx, cy := ebiten.CursorPosition()
		cursor := vector.New(float64(cx), float64(cy))
		if cursor.Within(vector.New(0, 0), f.screen) {
			target = cursor
		}
	}

	f.wg.Add(f.Conf.GoRoutines)
	for i := 0; i < f.Conf.GoRoutines; i++ {
		f.signal <- target
	}
	f.wg.Wait()
}

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type stats struct {
	Pos        vector.V
	Vel        vector.V
	Target     vector.V
	Cohesion   vector.V
	Separation vector.V
	Alignment  vector.V
	Targeting  vector.V
}

var leaderStats = stats{}

const neighbourRange float64 = 75

func (f *Flock) stepBoid(b *Boid, target vector.V) {
	if f.update {
		num := 0.0
		coh := vector.New(0, 0)
		sep := vector.New(0, 0)
		ali := vector.New(0, 0)
		for _, n := range f.Boids {
			if n == b {
				continue
			}
			diff := n.Pos.Subv(b.Pos)
			dist := diff.Length()
			if dist > neighbourRange {
				continue
			}
			num += 1
			coh = coh.Addv(n.Pos)
			sep = sep.Subv(separation(n, diff, dist))
			ali = ali.Addv(n.Vel)
		}
		if num > 0 {
			coh = cohesion(b, coh, num)
			ali = alignment(b, ali, num)
		}
		tar := centerTarget(b, target)
		b.Vel = b.Vel.Addv(coh).Addv(sep).Addv(ali).Addv(tar)
		b.Vel = clampSpeed(b)
		if b.ID == 0 {
			leaderStats = stats{b.Pos, b.Vel, target, coh, sep, ali, tar}
		}
	}
	b.Pos = b.Pos.Addv(b.Vel.Round())
}

const cohesionFactor float64 = 0.001

func cohesion(b *Boid, coh vector.V, num float64) vector.V {
	return coh.Div(num).Subv(b.Pos).Mul(cohesionFactor)
}

const separationRange float64 = 20
const separationFactor = 0.3

func separation(n *Boid, diff vector.V, dist float64) vector.V {
	if dist < separationRange {
		return diff.Mul((1 / dist) * separationFactor)
	}
	return vector.New(0, 0)
}

const alignmentFactor float64 = 0.1

func alignment(b *Boid, ali vector.V, num float64) vector.V {
	return ali.Div(num).Subv(b.Vel).Mul(alignmentFactor)
}

const targetRange float64 = 50
const targetRepelFactor float64 = 0.3
const targetAttractFactor float64 = 0.0001

func centerTarget(b *Boid, target vector.V) vector.V {
	diff := target.Subv(b.Pos)
	dist := diff.Length()
	if dist < targetRange {
		return diff.Mul((1 / dist) * -targetRepelFactor)
	}
	return diff.Mul(targetAttractFactor)
}

const velMax float64 = 1
const velMin float64 = 0.5

func clampSpeed(b *Boid) vector.V {
	l := b.Vel.Length()
	switch {
	case l > velMax:
		return b.Vel.Mul(velMax / l)
	case l < velMin:
		return b.Vel.Mul(velMin / l)
	}
	return b.Vel
}
