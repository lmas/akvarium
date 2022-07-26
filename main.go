package main

import (
	"flag"

	"github.com/lmas/boids/simulation"
)

var (
	flagDebug     = flag.Bool("debug", false, "Toggle debug info")
	flagEffects   = flag.Bool("effects", true, "Show extra graphic effects")
	flagInitSteps = flag.Int("initsteps", 0, "Run initial steps to prime the simulation")
)

func main() {
	flag.Parse()
	conf := simulation.DefaultConf()
	conf.Debug = *flagDebug
	conf.Effects = *flagEffects

	s, err := simulation.New(conf)
	if err != nil {
		panic(err)
	}

	if *flagInitSteps > 0 {
		s.Init(*flagInitSteps)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
