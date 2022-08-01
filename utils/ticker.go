package utils

import "math"

type Ticker struct {
	period float64
	rate   float64
	tick   float64
}

func NewTicker(ticksPerSecond, updatesPerSecond int) *Ticker {
	tps := float64(ticksPerSecond)
	ups := float64(updatesPerSecond)
	return &Ticker{
		period: (tps / ups) / tps,
		rate:   1.0 / tps,
		tick:   0,
	}
}

const tickerPrecision float64 = 1000

func (t *Ticker) round(f float64) float64 {
	return math.Round(f*tickerPrecision) / tickerPrecision
}

const reset float64 = 10000

func (t *Ticker) Tick() float64 {
	t.tick += t.rate
	if t.tick >= reset {
		t.tick = 0
	}
	return t.tick
}

func (t *Ticker) Float64() float64 {
	return t.tick
}

func (t *Ticker) Float32() float32 {
	return float32(t.tick)
}

func (t *Ticker) Mod(f float64) float64 {
	return math.Mod(t.round(t.tick/t.period), f)
}
