package boids

import (
	"math"
)

type IndexKey [2]int

// IndexBin contains IDs for Boids in the same bin.
type IndexBin []int

// Index groups Boids into neighbouring bins.
type Index struct {
	idx    indexMap
	offset float64
}

type indexMap map[IndexKey]IndexBin

func NewIndex(offset int) *Index {
	return &Index{
		idx:    make(indexMap),
		offset: float64(offset),
	}
}

// Key returns the key for the neighbouring bin a Boid is part of.
func (i *Index) Key(b *Boid) IndexKey {
	v := b.Pos.Div(i.offset)
	return IndexKey{
		int(math.Floor(v.X)),
		int(math.Floor(v.Y)),
	}
}

// Update clears the index and reinserts all Boids into new neighbouring bins.
func (i *Index) Update(boids []*Boid) {
	i.idx = make(indexMap)
	for _, b := range boids {
		k := i.Key(b)
		i.idx[k] = append(i.idx[k], b.ID)
	}
}

func (i *Index) IterBounds(min, max Vector, fun func(int)) {
	a, b, c, d := int(min.X), int(min.Y), int(max.X/i.offset), int(max.Y/i.offset)
	for k := range i.idx {
		if k[0] < a || k[1] < b || k[0] > c || k[1] > d {
			continue
		}
		for _, n := range i.idx[k] {
			fun(n)
		}
	}
}

// IterNeighbours iterates over all Boids in the same bin and the 8 neighbouring bins.
func (i *Index) IterNeighbours(b *Boid, fun func(n int)) {
	k := i.Key(b)
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			i.iterBin(IndexKey{k[0] + x, k[1] + y}, b.ID, fun)
		}
	}
}

func (i *Index) iterBin(k IndexKey, id int, fun func(n int)) {
	for _, n := range i.idx[k] {
		if n == id {
			continue
		}
		fun(n)
	}
}
