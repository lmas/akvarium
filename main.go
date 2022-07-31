package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lmas/boids/boids"
)

var (
	flagDebug   = flag.Bool("debug", false, "Toggle debug info")
	flagPretty  = flag.Bool("pretty", true, "Show pretty graphic effects")
	flagProfile = flag.Bool("profile", false, "Perform a CPU/MEM profile and quit")
	flagInit    = flag.Int("init", 2000, "Run initial updates to prime the simulation")
)

func main() {
	flag.Parse()
	if *flagProfile {
		go profileSim(".stats/cpu", ".stats/mem", 10)
	}

	conf := SimConf{
		Debug:         *flagDebug,
		Pretty:        *flagPretty,
		ScreenWidth:   1280,
		ScreenHeight:  720,
		UpdatesPerSec: 10,
		Swarm: boids.Conf{
			Boids:       500,
			Seed:        0,
			Workers:     10,
			IndexOffset: 50,
		},
	}

	s, err := New(conf)
	if err != nil {
		panic(err)
	}

	if *flagInit > 0 {
		s.Init(*flagInit)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SimConf struct {
	Debug         bool
	Pretty        bool
	ScreenWidth   int
	ScreenHeight  int
	UpdatesPerSec int
	Swarm         boids.Conf
}

type Simulation struct {
	Conf   SimConf
	boid   *ebiten.Image
	op     *ebiten.DrawImageOptions
	sop    *ebiten.DrawRectShaderOptions
	shader *ebiten.Shader
	swarm  *boids.Swarm
	screen boids.Vector
	target boids.Vector
	tick   *Ticker
}

//go:embed assets/boid.png
//go:embed assets/shader.go
var assets embed.FS

func New(conf SimConf) (*Simulation, error) {
	if conf.Swarm.Spawn[0].Length() == 0 && conf.Swarm.Spawn[1].Length() == 0 {
		conf.Swarm.Spawn = [2]boids.Vector{
			boids.NewVector(0, 0),
			boids.NewVector(float64(conf.ScreenWidth), float64(conf.ScreenHeight)),
		}
	}

	s := &Simulation{
		Conf:  conf,
		swarm: boids.New(conf.Swarm),
		op: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
		},
		sop: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]interface{}{
				"Resolution": []float32{
					float32(conf.ScreenWidth),
					float32(conf.ScreenHeight),
				},
			},
		},
		screen: boids.NewVector(float64(conf.ScreenWidth), float64(conf.ScreenHeight)),
		tick:   NewTicker(ebiten.MaxTPS(), conf.UpdatesPerSec),
	}
	s.Log("Loading assets..")

	sprite, err := loadImg("assets/boid.png")
	if err != nil {
		return nil, err
	}
	s.boid = ebiten.NewImageFromImage(sprite)

	b, err := assets.ReadFile("assets/shader.go")
	if err != nil {
		return nil, err
	}
	s.shader, err = ebiten.NewShader(b)
	if err != nil {
		return nil, err
	}

	s.Log("Assets ready")
	return s, nil
}

func (s *Simulation) Log(msg string, args ...interface{}) {
	if s.Conf.Debug {
		log.Printf(msg+"\n", args...)
	}
}

func (s *Simulation) Init(simulationSteps int) {
	s.Log("Priming simulation..")
	t := s.screen.Div(2)
	for i := 0; i < simulationSteps; i++ {
		// Must alternate between updating velocity (dirty) and position (non-dirty)
		s.swarm.Update(i%2 == 0, t)
	}
	s.Log("Simulation ready")
}

func (s *Simulation) Run() error {
	s.Log("Running simulation..")
	ebiten.SetWindowSize(s.Conf.ScreenWidth, s.Conf.ScreenHeight)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(s); err != nil {
		if !errors.Is(err, errQuit) {
			s.Log("Simulation error")
			return err
		}
	}
	s.Log("Simulation shutdown")
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// These 3 funcs are required by ebiten

func (s *Simulation) Layout(width, height int) (int, int) {
	return s.Conf.ScreenWidth, s.Conf.ScreenHeight
}

var errQuit = errors.New("quit")

func (s *Simulation) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return errQuit
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	s.tick.Tick()
	dirty := s.tick.Mod(1) == 0
	if dirty {
		cx, cy := ebiten.CursorPosition()
		cur := boids.NewVector(float64(cx), float64(cy))
		if cur.Within(minVec, s.screen) {
			s.target = cur
		} else {
			s.target = s.screen.Div(2)
		}
	}
	s.swarm.Update(dirty, s.target)
	return nil
}

// https://www.color-name.com/light-ocean-blue.color
var colBG = color.RGBA{0x04, 0x78, 0x9B, 0xFF}

// This prevents pop-in of boids at the top of the screen.
var minVec = boids.NewVector(-1, -1)

func (s *Simulation) Draw(screen *ebiten.Image) {
	screen.Fill(colBG)
	if s.Conf.Debug {
		s.drawDebug(screen)
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %0.f", ebiten.CurrentFPS()))
	}

	s.swarm.Index.IterBounds(minVec, s.screen, func(n int) {
		b := s.swarm.Boids[n]
		rotateAndTranslate(b.Pos, b.Vel.Angle(), s.boid, s.op)
		screen.DrawImage(s.boid, s.op)
		s.op.GeoM.Reset()
	})

	if s.Conf.Pretty {
		s.sop.Uniforms["Time"] = s.tick.Float32()
		screen.DrawRectShader(s.Conf.ScreenWidth, s.Conf.ScreenHeight, s.shader, s.sop)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// DEBUG

const rad2deg float64 = -180 / math.Pi

var colGreen = color.RGBA{0x0, 0xff, 0x0, 0x88}
var colRed = color.RGBA{0xff, 0x0, 0x0, 0x88}

func (s *Simulation) drawDebug(screen *ebiten.Image) {
	leader := s.swarm.Boids[0]
	k := s.swarm.Index.Key(leader)
	r := float64(s.Conf.Swarm.IndexOffset)

	// Shows bins around leader
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
		ebitenutil.DrawLine(screen, leader.Pos.X, leader.Pos.Y, n.Pos.X, n.Pos.Y, colGreen)
	})

	// Shows leader pos
	l := leader.Pos.Sub(r / 2)
	ebitenutil.DrawRect(screen, l.X, l.Y, r, r, colRed)

	// Shows target pos
	t := s.target.Sub(r / 2)
	ebitenutil.DrawRect(screen, t.X, t.Y, r, r, colRed)

	msg := fmt.Sprintf("TPS: %0.f  FPS: %0.f  Target: %0.f,%0.f  Leader: %3.0f,%3.0f  %s  %+0.1f°\n",
		ebiten.CurrentTPS(), ebiten.CurrentFPS(),
		s.target.X, s.target.Y,
		leader.Pos.X, leader.Pos.Y,
		leader.Vel, leader.Vel.Angle()*rad2deg,
	)
	ebitenutil.DebugPrint(screen, msg)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// UTILS

func profileSim(cpu, mem string, sleep int) {
	c, err := os.Create(cpu)
	if err != nil {
		panic(err)
	}
	m, err := os.Create(mem)
	if err != nil {
		panic(err)
	}

	defer func() {
		c.Close()
		runtime.GC()
		pprof.WriteHeapProfile(m)
		m.Close()
		os.Exit(0)
	}()

	pprof.StartCPUProfile(c)
	time.Sleep(time.Duration(sleep) * time.Second)
	pprof.StopCPUProfile()
}

func loadImg(p string) (image.Image, error) {
	f, err := assets.Open(p)
	if err != nil {
		return nil, fmt.Errorf("could not open '%s': %s", p, err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode '%s': %s", p, err)
	}
	return i, nil
}

const maxAngleN float64 = -math.Pi / 2
const maxAngleNE float64 = -math.Pi / 4
const maxAngleSE float64 = math.Pi / 4
const maxAngleS float64 = math.Pi / 2

func clampAngleAndFlip(a float64) (float64, bool) {
	flipped := false
	if a < maxAngleN {
		a += math.Pi
		flipped = true
	} else if a > maxAngleS {
		a -= math.Pi
		flipped = true
	}
	if a > maxAngleN && a < maxAngleNE {
		a = maxAngleNE
	} else if a < maxAngleS && a > maxAngleSE {
		a = maxAngleSE
	}
	return a, flipped
}

func rotateAndTranslate(pos boids.Vector, angle float64, src *ebiten.Image, op *ebiten.DrawImageOptions) {
	x, y := src.Size()
	w, h := float64(x), float64(y)
	a, flipped := clampAngleAndFlip(angle)
	if flipped {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(w, 0)
	}
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(a)
	op.GeoM.Translate(pos.X, pos.Y)
}

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
