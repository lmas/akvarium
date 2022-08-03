package main

// This tool let's you try out shaders quick 'n dirty.

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lmas/boids/utils"
)

var (
	flagShader = flag.String("p", "", "Path to shader to run")
)

const screenWidth int = 1280
const screenHeight int = 720

var errQuit = errors.New("quit")

type shaderSim struct {
	sprite *ebiten.Image
	op     *ebiten.DrawImageOptions
	shader *ebiten.Shader
	sop    *ebiten.DrawRectShaderOptions
	tick   *utils.Ticker
	dir    float64
	pos    float64
}

func main() {
	flag.Parse()

	sim := &shaderSim{
		sprite: loadSprite("assets/boid-gopher.png"),
		op: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
		},
		shader: loadShader(*flagShader),
		sop: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]interface{}{
				"Resolution": []float32{
					float32(screenWidth),
					float32(screenHeight),
				},
			},
		},
		tick: utils.NewTicker(ebiten.MaxTPS(), 10),
	}
	w, h := sim.sprite.Size()
	sim.op.GeoM.Translate(float64(screenWidth)/2, 0)
	sim.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shader Test")
	if err := ebiten.RunGame(sim); err != nil {
		if !errors.Is(err, errQuit) {
			panic(err)
		}
	}
}

func (s *shaderSim) Layout(width, height int) (int, int) {
	return screenWidth, screenHeight
}

func (s *shaderSim) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return errQuit
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if s.dir == 0 {
		s.dir = 1
	}
	s.pos += s.dir
	if s.pos <= 0 {
		s.dir = 1
	} else if s.pos >= float64(screenHeight) {
		s.dir = -1
	}
	s.tick.Tick()
	return nil
}

var colBG = color.RGBA{0x04, 0x78, 0x9B, 0xFF}

func (s *shaderSim) Draw(screen *ebiten.Image) {
	screen.Fill(colBG)
	s.op.GeoM.Translate(0, s.dir)
	screen.DrawImage(s.sprite, s.op)
	s.sop.Uniforms["Time"] = s.tick.Float32()
	screen.DrawRectShader(screenWidth, screenHeight, s.shader, s.sop)
	msg := fmt.Sprintf("TPS: %0.1f  FPS: %0.1f  Tick: %0.1f", ebiten.CurrentTPS(), ebiten.CurrentFPS(), s.tick.Float64())
	ebitenutil.DebugPrint(screen, msg)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func loadSprite(p string) *ebiten.Image {
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(i)
}

func loadShader(p string) *ebiten.Shader {
	b, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}
	shader, err := ebiten.NewShader(b)
	if err != nil {
		panic(err)
	}
	return shader
}
