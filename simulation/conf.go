package simulation

type Conf struct {
	Debug               bool
	Seed                int64
	GoRoutines          int
	ScreenWidth         int
	ScreenHeight        int
	FlockSize           int
	SpeedLimitingFactor float64
	VisionRadious       float64
	SeparationRadious   float64
	ScreenScale         float64
	SeparationFactor    float64
	CohesionFactor      float64
	AlignmentFactor     float64
	TargetingFactor     float64
}

func DefaultConf() Conf {
	return Conf{
		Debug:               true,
		Seed:                0,
		GoRoutines:          10,
		ScreenWidth:         1280,
		ScreenHeight:        720,
		FlockSize:           500,
		SpeedLimitingFactor: 1,
		VisionRadious:       50,
		SeparationRadious:   20,
		ScreenScale:         0.1,
		SeparationFactor:    0.003,
		CohesionFactor:      0.00001,
		AlignmentFactor:     0.001,
		TargetingFactor:     0.00002,
	}
}
