package main

import (
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
	"github.com/lmas/boids/assets"
	"github.com/lmas/boids/boids"
	"github.com/lmas/boids/vector"
)

var (
	flagDebug   = flag.Bool("debug", false, "Toggle debug info")
	flagEffects = flag.Bool("effects", true, "Show extra graphic effects")
	flagInit    = flag.Int("init", 2000, "Run initial update steps to prime the simulation")
)

func main() {
	flag.Parse()
	conf := SimConf{
		Debug:        *flagDebug,
		Effects:      *flagEffects,
		ScreenWidth:  1280,
		ScreenHeight: 720,
		Swarm: boids.Conf{
			Seed:        0,
			GoRoutines:  10,
			SwarmSize:   500,
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

var (
	colGreen = color.RGBA{0x0, 0xff, 0x0, 0x88}
	colRed   = color.RGBA{0xff, 0x0, 0x0, 0x88}
	errQuit  = errors.New("quit")
)

type SimConf struct {
	Debug        bool
	Effects      bool
	ScreenWidth  int
	ScreenHeight int
	Swarm        boids.Conf
}

type Simulation struct {
	Conf     SimConf
	boidSize [2]float64
	boidImg  *ebiten.Image
	imgOP    *ebiten.DrawImageOptions
	swarm    *boids.Swarm
	maxTPS   int
	tps      int
	tick     int
	screen   vector.V
	target   vector.V
}

const screenScale float64 = 0.04 // Scales down the sprite

func New(conf SimConf) (*Simulation, error) {
	s := &Simulation{
		Conf:   conf,
		swarm:  boids.New(conf.Swarm),
		imgOP:  &ebiten.DrawImageOptions{},
		screen: vector.New(float64(conf.ScreenWidth), float64(conf.ScreenHeight)),
		maxTPS: ebiten.MaxTPS(),
	}
	s.Log("Loading assets..")

	f, err := assets.FS.Open("shiny_boid.png")
	if err != nil {
		return s, fmt.Errorf("Failed to open boid sprite: %s", err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return s, fmt.Errorf("Failed to decode boid sprite: %s", err)
	}

	img := ebiten.NewImageFromImage(i)
	w, h := img.Size()
	s.boidSize[0], s.boidSize[1] = float64(w)*screenScale, float64(h)*screenScale
	s.boidImg = ebiten.NewImage(int(s.boidSize[0]), int(s.boidSize[1]))
	s.imgOP.GeoM.Scale(screenScale, screenScale)
	s.boidImg.DrawImage(img, s.imgOP)

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
		// Must alternate between dirty (updates velocity, expensive) and non-dirty (updates position, cheap).
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

const tickLimiter int = 6

var zeroVec = vector.New(0, 0)

func (s *Simulation) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return errQuit
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	s.tick += 1
	if s.tick >= s.maxTPS {
		s.tick = 0
	}
	dirty := false
	if s.tick%tickLimiter == 0 {
		dirty = true
		cx, cy := ebiten.CursorPosition()
		cur := vector.New(float64(cx), float64(cy))
		if cur.Within(zeroVec, s.screen) {
			s.target = cur
		} else {
			s.target = s.screen.Div(2)
		}
	}
	s.swarm.Update(dirty, s.target)
	return nil
}

const shiftAngle float64 = math.Pi / 2 // Shifts the sprite by 90 degrees

func (s *Simulation) Draw(screen *ebiten.Image) {
	if s.Conf.Debug {
		s.drawDebug(screen)
	}
	for _, b := range s.swarm.Boids {
		if s.Conf.Effects {
			s.imgOP.ColorM.Reset()
			hue := b.Vel.Angle() * 0.05
			scale := (b.Pos.Angle() + hue)
			s.imgOP.ColorM.ChangeHSV(hue, 1, scale)
		}

		s.imgOP.GeoM.Reset()
		s.imgOP.GeoM.Translate(-s.boidSize[0]/2, -s.boidSize[1]/2)
		s.imgOP.GeoM.Rotate(b.Vel.Angle() + shiftAngle)
		s.imgOP.GeoM.Translate(b.Pos.X, b.Pos.Y)
		screen.DrawImage(s.boidImg, s.imgOP)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const rad2deg float64 = -180 / math.Pi

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

	// Shows all bins
	s.swarm.Index.IterBins(func(bin boids.IndexKey) {
		x, y := float64(bin[0])*r, float64(bin[1])*r
		ebitenutil.DrawRect(screen, x, y, r, r, colGreen)
		ebitenutil.DrawLine(screen, x, y, x+r, y, colGreen)
		ebitenutil.DrawLine(screen, x, y, x, y+r, colGreen)
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
