package simulation

type Conf struct {
	Debug        bool
	Seed         int64
	GoRoutines   int
	SwarmSize    int
	ScreenWidth  int
	ScreenHeight int
	ScreenScale  float64
}

func DefaultConf() Conf {
	return Conf{
		Debug:        true,
		Seed:         0,
		GoRoutines:   10,
		SwarmSize:    500,
		ScreenWidth:  1280,
		ScreenHeight: 720,
		ScreenScale:  0.1,
	}
}
