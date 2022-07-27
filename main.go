package main

import (
	"flag"

	"github.com/lmas/boids/simulation"
)

var (
	flagDebug   = flag.Bool("debug", false, "Toggle debug info")
	flagEffects = flag.Bool("effects", true, "Show extra graphic effects")
	flagInit    = flag.Int("init", 2000, "Run initial update steps to prime the simulation")
)

func main() {
	flag.Parse()
	conf := simulation.SimConf{
		Debug:        *flagDebug,
		Effects:      *flagEffects,
		ScreenWidth:  1280,
		ScreenHeight: 720,
		Swarm: simulation.Conf{
			Seed:       0,
			GoRoutines: 10,
			SwarmSize:  500,
		},
	}

	s, err := simulation.New(conf)
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
