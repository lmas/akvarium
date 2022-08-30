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
		sep = sep.Subv(s.separation(b, n))
	})

	if num > 0 {
		coh = s.cohesion(b, coh, num)
		ali = s.alignment(b, ali, num)
	}
	tar := s.centerTarget(b, target)
	b.Vel = b.Vel.Addv(coh).Addv(ali).Addv(sep).Addv(tar)
	b.Vel = s.clampSpeed(b)
}

func (s *Swarm) cohesion(b *Boid, coh Vector, num float64) Vector {
	return coh.Div(num).Subv(b.Pos).Mul(s.Conf.CohesionFactor)
}

func (s *Swarm) alignment(b *Boid, ali Vector, num float64) Vector {
	return ali.Div(num).Subv(b.Vel).Mul(s.Conf.AlignmentFactor)
}

func (s *Swarm) separation(b, n *Boid) Vector {
	diff := n.Pos.Subv(b.Pos)
	dist := diff.InRange(s.squareSeparationRange)
	if dist > 0 {
		return diff.Div(dist / s.Conf.SeparationFactor)
	}
	return NewVector(0, 0)
}

func (s *Swarm) centerTarget(b *Boid, target Vector) Vector {
	diff := target.Subv(b.Pos)
	dist := diff.InRange(s.squareTargetRange)
	if dist > 0 {
		return diff.Div(dist / -s.Conf.TargetRepelFactor)
	}
	return diff.Mul(s.Conf.TargetAttractFactor)
}

func (s *Swarm) clampSpeed(b *Boid) Vector {
	l := b.Vel.Dot(b.Vel)
	switch {
	case l > s.squareVelocityMax:
		return b.Vel.Mul(s.Conf.VelocityMax / math.Sqrt(l))
	case l < s.squareVelocityMin:
		return b.Vel.Mul(s.Conf.VelocityMin / math.Sqrt(l))
	}
	return b.Vel
}
