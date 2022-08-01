package utils

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func RunProfiler(cpu, mem string, sleep int) {
	c, err := os.Create(cpu)
	if err != nil {
		panic(err)
	}
	m, err := os.Create(mem)
	if err != nil {
		panic(err)
	}

	defer func() {
		c.Close()
		runtime.GC()
		err = pprof.WriteHeapProfile(m)
		if err != nil {
			fmt.Println(err)
		}
		m.Close()
		os.Exit(0)
	}()

	err = pprof.StartCPUProfile(c)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(sleep) * time.Second)
	pprof.StopCPUProfile()
}
