package boids

import (
	"math"
)

type IndexKey [2]int

type IndexBin []int

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

func (i *Index) Key(b *Boid) IndexKey {
	v := b.Pos.Div(i.offset)
	return IndexKey{
		int(math.Floor(v.X)),
		int(math.Floor(v.Y)),
	}
}

func (i *Index) Update(boids []*Boid) {
	i.idx = make(indexMap)
	for _, b := range boids {
		k := i.Key(b)
		i.idx[k] = append(i.idx[k], b.ID)
	}
}

func (i *Index) IterBins(fun func(IndexKey)) {
	for b := range i.idx {
		fun(b)
	}
}

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
