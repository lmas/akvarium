package simulation

import (
	"errors"
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
)

var (
	colGreen = color.RGBA{0x0, 0xff, 0x0, 0x88}
	colRed   = color.RGBA{0xff, 0x0, 0x0, 0x88}
	errQuit  = errors.New("quit")
)

type Simulation struct {
	Conf     Conf
	boidSize [2]float64
	boidImg  *ebiten.Image
	imgOP    *ebiten.DrawImageOptions
	swarm    *Swarm
	maxTPS   int
	tps      int
	tick     int
}

const screenScale float64 = 0.04 // Scales down the sprite

func New(conf Conf) (*Simulation, error) {
	s := &Simulation{
		Conf: conf,
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

	s.imgOP = &ebiten.DrawImageOptions{}
	s.imgOP.GeoM.Scale(screenScale, screenScale)
	img := ebiten.NewImageFromImage(i)
	w, h := img.Size()
	s.boidSize[0], s.boidSize[1] = float64(w)*screenScale, float64(h)*screenScale
	s.boidImg = ebiten.NewImage(int(s.boidSize[0]), int(s.boidSize[1]))
	s.boidImg.DrawImage(img, s.imgOP)

	s.maxTPS = ebiten.MaxTPS()
	s.swarm = NewSwarm(conf)

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
	s.swarm.Init(simulationSteps)
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
	s.swarm.Step(s.tick%tickLimiter == 0) // Limit the amount of dirty Steps(). If TPS=60, updates=6
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
	k := getBinKey(leader)
	nr := neighbourRange
	sr := separationRange
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			x := float64(k[0]+i) * nr
			y := float64(k[1]+j) * nr
			ebitenutil.DrawRect(screen, x, y, nr, nr, colGreen)
			ebitenutil.DrawLine(screen, x, y, x+nr, y, colGreen)
			ebitenutil.DrawLine(screen, x, y, x, y+nr, colGreen)
		}
	}
	for p := range s.swarm.index {
		x, y := float64(p[0])*nr, float64(p[1])*nr
		ebitenutil.DrawRect(screen, x, y, nr, nr, colGreen)
		ebitenutil.DrawLine(screen, x, y, x+nr, y, colGreen)
		ebitenutil.DrawLine(screen, x, y, x, y+nr, colGreen)
	}
	ebitenutil.DrawRect(screen, leader.Pos.X-sr, leader.Pos.Y-sr, sr*2, sr*2, colRed)
	t := leaderStats.Target.Sub(targetRange / 2)
	ebitenutil.DrawRect(screen, t.X, t.Y, targetRange, targetRange, colRed)

	msg := fmt.Sprintf("TPS: %0.f  FPS: %0.f  Target: %0.f,%0.f  Leader: %3.0f,%3.0f  %s  %+0.1f°\n"+
		"coh: %s  sep: %s  ali: %s  tar: %s",
		ebiten.CurrentTPS(), ebiten.CurrentFPS(),
		leaderStats.Target.X, leaderStats.Target.Y,
		leaderStats.Pos.X, leaderStats.Pos.Y,
		leaderStats.Vel, leaderStats.Vel.Angle()*rad2deg,
		leaderStats.Cohesion, leaderStats.Separation,
		leaderStats.Alignment, leaderStats.Targeting,
	)
	ebitenutil.DebugPrint(screen, msg)
}
