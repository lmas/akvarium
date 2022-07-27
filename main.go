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
	conf := simulation.DefaultConf()
	conf.Debug = *flagDebug
	conf.Effects = *flagEffects

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
