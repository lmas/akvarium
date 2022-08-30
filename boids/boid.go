package boids

import "math"

// Boid represents a single boid.
// It will try to fit in with a Swarm by:
// - moving towards the center of nearby Boids (Cohesion).
// - matching with nearby Boids' velocity (Alignment).
// - avoiding collisions with nearby Boids (Separation).
//
// It can optionally move towards a target.
type Boid struct {
	ID  int
	Pos Vector
	Vel Vector
}

func (s *Swarm) updateBoid(b *Boid, dirty bool, target Vector) {
	if !dirty {
		b.Pos = b.Pos.Addv(b.Vel.Round())
		return
	}

	num := 0.0
	coh := NewVector(0, 0)
	ali := NewVector(0, 0)
	sep := NewVector(0, 0)
	s.Index.IterNeighbours(b, func(id int) {
		n := s.Boids[id]
		num += 1
		coh = coh.Addv(n.Pos)
		ali = ali.Addv(n.Vel)
		sep = sep.Subv(separation(b, n))
	})

	if num > 0 {
		coh = cohesion(b, coh, num)
		ali = alignment(b, ali, num)
	}
	tar := centerTarget(b, target)
	b.Vel = b.Vel.Addv(coh).Addv(ali).Addv(sep).Addv(tar)
	b.Vel = clampSpeed(b)
}

const cohesionFactor float64 = 0.001

func cohesion(b *Boid, coh Vector, num float64) Vector {
	return coh.Div(num).Subv(b.Pos).Mul(cohesionFactor)
}

const alignmentFactor float64 = 0.05

func alignment(b *Boid, ali Vector, num float64) Vector {
	return ali.Div(num).Subv(b.Vel).Mul(alignmentFactor)
}

const separationRange float64 = 20
const separationFactor = 0.3

const tsep float64 = separationRange * separationRange

func separation(b, n *Boid) Vector {
	diff := n.Pos.Subv(b.Pos)
	dist := diff.InRange(tsep)
	if dist > 0 {
		return diff.Div(dist / separationFactor)
	}
	return NewVector(0, 0)
}

const targetRange float64 = 50
const targetRepelFactor float64 = 0.3
const targetAttractFactor float64 = 0.00004

const ttar float64 = targetRange * targetRange

func centerTarget(b *Boid, target Vector) Vector {
	diff := target.Subv(b.Pos)
	dist := diff.InRange(ttar)
	if dist > 0 {
		return diff.Div(dist / -targetRepelFactor)
	}
	return diff.Mul(targetAttractFactor)
}

const velMax float64 = 1
const velMin float64 = 0.5

const tvmax float64 = velMax * velMax
const tvmin float64 = velMin * velMin

func clampSpeed(b *Boid) Vector {
	l := b.Vel.Dot(b.Vel)
	switch {
	case l > tvmax:
		return b.Vel.Mul(velMax / math.Sqrt(l))
	case l < tvmin:
		return b.Vel.Mul(velMin / math.Sqrt(l))
	}
	return b.Vel
}
