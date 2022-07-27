package boids

import (
	"github.com/lmas/boids/vector"
)

type Boid struct {
	ID  int
	Pos vector.V
	Vel vector.V
}

const neighbourRange float64 = 50

func (s *Swarm) updateBoid(b *Boid, dirty bool, target vector.V) {
	if !dirty {
		b.Pos = b.Pos.Addv(b.Vel.Round())
		return
	}

	num := 0.0
	coh := vector.New(0, 0)
	sep := vector.New(0, 0)
	ali := vector.New(0, 0)
	s.IterNeighbours(b, func(n *Boid) {
		num += 1
		coh = coh.Addv(n.Pos)
		sep = sep.Subv(separation(b, n))
		ali = ali.Addv(n.Vel)
	})

	if num > 0 {
		coh = cohesion(b, coh, num)
		ali = alignment(b, ali, num)
	}
	tar := centerTarget(b, target)
	b.Vel = b.Vel.Addv(coh).Addv(sep).Addv(ali).Addv(tar)
	b.Vel = clampSpeed(b)
}

const cohesionFactor float64 = 0.001

func cohesion(b *Boid, coh vector.V, num float64) vector.V {
	return coh.Div(num).Subv(b.Pos).Mul(cohesionFactor)
}

const separationRange float64 = 20
const separationFactor = 0.3

func separation(b, n *Boid) vector.V {
	diff := n.Pos.Subv(b.Pos)
	dist := diff.Length()
	if dist < separationRange {
		return diff.Div(dist / separationFactor)
	}
	return vector.New(0, 0)
}

const alignmentFactor float64 = 0.05

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
		return diff.Div(dist / -targetRepelFactor)
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
