package simulation

import (
	"math/rand"

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
}

func NewFlock(conf Conf) *Flock {
	f := &Flock{
		Boids:        make([]*Boid, conf.FlockSize),
		Conf:         conf,
		screenWidth:  float64(conf.ScreenWidth),
		screenHeight: float64(conf.ScreenHeight),
	}

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
	return f
}

func (f *Flock) Init(simulationSteps int) {
	for i := 0; i < simulationSteps; i++ {
		f.Step()
	}
}

func (f *Flock) Step() {
	sum := 0.0
	vel := vector.New(0, 0)
	cohesion := vector.New(0, 0)
	alignment := vector.New(0, 0)
	separation := vector.New(0, 0)
	targetpoint := vector.New(f.screenWidth/2, f.screenHeight/2)

	for _, current := range f.Boids {
		sum = 0.0
		vel = vel.Zero()
		cohesion = cohesion.Zero()
		alignment = alignment.Zero()
		separation = separation.Zero()

		for _, target := range f.Boids {
			if target == current || current.Pos.Distance(target.Pos) > f.Conf.VisionRadious {
				continue
			}
			sum += 1.0
			cohesion = cohesion.Addv(target.Pos)
			alignment = alignment.Addv(target.Vel)
			if current.Pos.Distance(target.Pos) < f.Conf.SeparationRadious {
				separation = separation.Addv(current.Pos.Subv(target.Pos))
				if separation.Length() < 1 {
					separation = vector.New(rand.Float64(), rand.Float64())
				}
			}
		}
		if sum > 0 {
			cohesion = cohesion.Div(sum).Subv(current.Pos)
			alignment = alignment.Div(sum).Subv(current.Vel)
		}

		vel = vel.Addv(cohesion.Mul(f.Conf.CohesionFactor))
		vel = vel.Addv(alignment.Mul(f.Conf.AlignmentFactor))
		vel = vel.Addv(separation.Mul(f.Conf.SeparationFactor))
		vel = vel.Addv(targetpoint.Subv(current.Pos).Mul(f.Conf.TargetingFactor))
		current.Vel = f.SpeedLimit(current.Vel.Addv(vel)).Round()
		current.Pos = current.Pos.Addv(current.Vel).Round()
	}
}

func (f *Flock) SpeedLimit(v vector.V) vector.V {
	len := v.Length()
	if len > f.Conf.SpeedLimitingFactor {
		v = v.Div(len)
	}
	return v.Mul(f.Conf.SpeedLimitingFactor)
}
