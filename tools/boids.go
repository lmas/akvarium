package main

// This tool let's you try out the boid simulation quickly.

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lmas/akvarium/boids"
	"github.com/lmas/akvarium/utils"
)

const screenWidth int = 1280
const screenHeight int = 720

var errQuit = errors.New("quit")
var minVec = boids.NewVector(-1, -1)
var maxVec = boids.NewVector(float64(screenWidth), float64(screenHeight))

var conf = boids.Conf{
	Boids:       500,
	Seed:        0,
	Workers:     10,
	IndexOffset: 50,
	Spawn:       [2]boids.Vector{minVec, maxVec},
}

type debugSim struct {
	swarm  *boids.Swarm
	sprite *ebiten.Image
	op     *ebiten.DrawImageOptions
	target boids.Vector
	tick   *utils.Ticker
}

func main() {
	s := &debugSim{
		swarm: boids.New(conf),
		op: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
		},
		tick: utils.NewTicker(ebiten.MaxTPS(), 10),
	}

	f, err := os.Open("assets/boid-clownfish.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	s.sprite = ebiten.NewImageFromImage(i)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids Debug")
	if err := ebiten.RunGame(s); err != nil {
		if !errors.Is(err, errQuit) {
			panic(err)
		}
	}
}

func (s *debugSim) Layout(width, height int) (int, int) {
	return screenWidth, screenHeight
}

func (s *debugSim) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return errQuit
	}
	s.tick.Tick()
	dirty := s.tick.Mod(1) == 0
	if dirty {
		cx, cy := ebiten.CursorPosition()
		cur := boids.NewVector(float64(cx), float64(cy))
		if cur.Within(minVec, maxVec) {
			s.target = cur
		} else {
			s.target = maxVec.Div(2)
		}
	}
	s.swarm.Update(dirty, s.target)
	return nil
}

var colGreen = color.RGBA{0x0, 0xff, 0x0, 0x88}
var colRed = color.RGBA{0xff, 0x0, 0x0, 0x88}

func (s *debugSim) Draw(screen *ebiten.Image) {
	leader := s.swarm.Boids[0]
	// Shows bins around leader
	k := s.swarm.Index.Key(leader)
	r := float64(conf.IndexOffset)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			x := float64(k[0]+i) * r
			y := float64(k[1]+j) * r
			ebitenutil.DrawRect(screen, x, y, r, r, colGreen)
			ebitenutil.DrawLine(screen, x, y, x+r, y, colGreen)
			ebitenutil.DrawLine(screen, x, y, x, y+r, colGreen)
		}
	}

	// Show lines connecting leader with it's neighbours
	s.swarm.Index.IterNeighbours(leader, func(id int) {
		n := s.swarm.Boids[id]
		ebitenutil.DrawLine(screen, leader.Pos.X, leader.Pos.Y, n.Pos.X, n.Pos.Y, colRed)
	})

	// Shows target pos
	t := s.target.Sub(r / 2)
	ebitenutil.DrawRect(screen, t.X, t.Y, r, r, colRed)

	// Draw the boids
	x, y := s.sprite.Size()
	w, h := float64(x), float64(y)
	s.swarm.Index.IterBounds(minVec, maxVec, func(n int) {
		b := s.swarm.Boids[n]
		s.op.GeoM.Translate(-w/2, -h/2)
		s.op.GeoM.Rotate(b.Vel.Angle())
		s.op.GeoM.Translate(b.Pos.X, b.Pos.Y)
		screen.DrawImage(s.sprite, s.op)
		s.op.GeoM.Reset()
	})

	msg := fmt.Sprintf("TPS: %0.f  FPS: %0.f  Tick: %0.1f  Target: %0.f,%0.f  Leader: %3.0f,%3.0f  %s  %+0.1f°\n",
		ebiten.CurrentTPS(), ebiten.CurrentFPS(), s.tick.Float64(),
		s.target.X, s.target.Y,
		leader.Pos.X, leader.Pos.Y,
		leader.Vel, leader.Vel.Angle(),
	)
	ebitenutil.DebugPrint(screen, msg)
}
