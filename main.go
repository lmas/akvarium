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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lmas/akvarium/boids"
	"github.com/lmas/akvarium/utils"
)

const Version string = "v0.1.0"
const Banner string = `
                               #xo;:::x
                  #$xxxxx$##x;;::::::::o$$#
                $;::::::::    o:::::::::::::,;$
               #:::::::::o     o;;xxx;:::;;   ,:o$
 $:o$xxooxoxx$xx::::::::::.     o#o;:;:::,   x;$o::x
 ,o#::::::,    o::::::::::;,   ,$::::::::,  .:::::::o
 :;#:::::::    ::::::::::::,   .$::::::::x  ,::::::;$
 $.$o:::::: ,:;x:::::::::::    ooxo::::::::  :o;;o$
  $;o$ooxx#     o::::::::;$.,::$o;;::::::::;;o$
                #xo;;:;ox#      #::::::x#    Akvarium
                                #;;;oo#      `

var (
	flagInit    = flag.Int("init", 2000, "Run initial updates to prime the simulation")
	flagProfile = flag.Bool("profile", false, "Perform a CPU/MEM profile and exit after 30 seconds")
	flagVerbose = flag.Bool("verbose", false, "Toggle verbose info")
	flagVersion = flag.Bool("version", false, "Print version and exit")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf(Banner)
		fmt.Println(Version)
		return
	}

	conf := SimConf{
		Verbose:       *flagVerbose,
		ScreenWidth:   1280,
		ScreenHeight:  720,
		UpdatesPerSec: 10,
		Swarm: boids.Conf{
			Seed:                0,
			Boids:               500,
			Workers:             10,
			IndexOffset:         50,
			CohesionFactor:      0.001,
			AlignmentFactor:     0.05,
			SeparationRange:     20,
			SeparationFactor:    0.3,
			TargetRange:         50,
			TargetRepelFactor:   0.3,
			TargetAttractFactor: 0.00004,
			VelocityMax:         1,
			VelocityMin:         0.5,
		},
	}

	if *flagProfile {
		go utils.RunProfiler(".stats/cpu", ".stats/mem", 30)
	}

	s, err := New(conf)
	if err != nil {
		panic(err)
	}

	if !*flagProfile && *flagInit > 0 {
		s.Init(*flagInit)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SimConf struct {
	Verbose       bool
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
	tick   *utils.Ticker
}

//go:embed assets/boid-clownfish.png
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
		tick:   utils.NewTicker(ebiten.MaxTPS(), conf.UpdatesPerSec),
	}
	s.Log("Loading assets..")

	sprite, err := loadImg("assets/boid-clownfish.png")
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
	return s, nil
}

func (s *Simulation) Log(msg string, args ...interface{}) {
	if s.Conf.Verbose {
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
	if s.Conf.Verbose {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %0.f", ebiten.CurrentFPS()))
	}

	s.swarm.Index.IterBounds(minVec, s.screen, func(n int) {
		b := s.swarm.Boids[n]
		rotateAndTranslate(b.Pos, b.Vel.Angle(), s.boid, s.op)
		screen.DrawImage(s.boid, s.op)
		s.op.GeoM.Reset()
	})

	s.sop.Uniforms["Time"] = s.tick.Float32()
	screen.DrawRectShader(s.Conf.ScreenWidth, s.Conf.ScreenHeight, s.shader, s.sop)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// UTILS

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
