package simulation

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lmas/boids/assets"
)

var (
	colGreen = color.RGBA{0x0, 0xff, 0x0, 0x88}
	colRed   = color.RGBA{0xff, 0x0, 0x0, 0x88}

	errQuit = errors.New("quit")
)

type Simulation struct {
	Conf     Conf
	boidSize [2]float64
	boidImg  *ebiten.Image
	imgOP    *ebiten.DrawImageOptions
	flock    *Flock
	maxTPS   int
	tps      int
	tick     int
}

func New(conf Conf) (*Simulation, error) {
	s := &Simulation{
		Conf: conf,
	}
	s.Log("Loading assets..")

	f, err := assets.FS.Open("fishy.png")
	if err != nil {
		return s, fmt.Errorf("Failed to open boid sprite: %s", err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return s, fmt.Errorf("Failed to decode boid sprite: %s", err)
	}

	s.imgOP = &ebiten.DrawImageOptions{}
	s.imgOP.GeoM.Scale(conf.ScreenScale, conf.ScreenScale)
	img := ebiten.NewImageFromImage(i)
	w, h := img.Size()

	s.boidSize[0], s.boidSize[1] = float64(w)*conf.ScreenScale, float64(h)*conf.ScreenScale
	s.boidImg = ebiten.NewImage(int(s.boidSize[0]), int(s.boidSize[1]))
	s.boidImg.DrawImage(img, s.imgOP)
	s.flock = NewFlock(conf)
	s.maxTPS = ebiten.MaxTPS()
	s.tps = s.maxTPS / 10 // 10 times per second

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
	s.flock.Init(simulationSteps)
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

func (s *Simulation) Draw(screen *ebiten.Image) {
	for _, b := range s.flock.Boids {
		s.imgOP.GeoM.Reset()
		s.imgOP.GeoM.Translate(-s.boidSize[0]/2, -s.boidSize[1]/2)
		s.imgOP.GeoM.Rotate(b.Vel.Angle())
		s.imgOP.GeoM.Translate(b.Pos.X, b.Pos.Y)
		if s.Conf.Debug && b == s.flock.Boids[0] {
			vr := float64(neighbourRange)
			ebitenutil.DrawRect(screen, b.Pos.X-vr, b.Pos.Y-vr, vr*2, vr*2, colGreen)
			sr := float64(separationRange)
			ebitenutil.DrawRect(screen, b.Pos.X-sr, b.Pos.Y-sr, sr*2, sr*2, colRed)
		}
		screen.DrawImage(s.boidImg, s.imgOP)
	}

	if s.Conf.Debug {
		t := leaderStats.Target.Sub(targetRange / 2)
		ebitenutil.DrawRect(screen, t.X, t.Y, targetRange, targetRange, colRed)
		msg := fmt.Sprintf("TPS: %0.f  FPS: %0.f  Target: %0.f,%0.f  Leader: %3.0f,%3.0f  %s  %+0.1f°\n"+
			"coh: %s  sep: %s  ali: %s  tar: %s",
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
			leaderStats.Target.X, leaderStats.Target.Y,
			leaderStats.Pos.X, leaderStats.Pos.Y,
			leaderStats.Vel,
			leaderStats.Vel.Angle(),
			leaderStats.Cohesion,
			leaderStats.Separation,
			leaderStats.Alignment,
			leaderStats.Targeting,
		)
		ebitenutil.DebugPrint(screen, msg)
	}
}

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
	if s.tick%s.tps == 0 {
		s.flock.Step(true)
	}
	s.flock.Step(false)
	return nil
}
