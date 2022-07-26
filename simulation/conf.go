package simulation

type Conf struct {
	Debug        bool
	Effects      bool
	Seed         int64
	GoRoutines   int
	SwarmSize    int
	ScreenWidth  int
	ScreenHeight int
}

func DefaultConf() Conf {
	return Conf{
		Debug:        true,
		Effects:      true,
		Seed:         0,
		GoRoutines:   10,
		SwarmSize:    500,
		ScreenWidth:  1280,
		ScreenHeight: 720,
	}
}
