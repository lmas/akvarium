
# Akvarium

A tiny aquarium running a [Boids] simulation.

![banner]

> Boids is an artificial life program,
> developed by [Craig Reynolds] in 1986,
> which simulates the flocking behaviour of birds.
> - Wikipedia

Using only three simple rules one is able to simulate emergent flocking behaviour that is mesmerising.
The three rules are:

(pic of rules)

- **Cohesion:** A single Boid tries to move towards the center of a nearby group of other Boids.
- **Alignment:** And it should try to match it's velocity and direction with it's neighbours.
- **Separation:** While moving, it should also try to avoid collisions with the closest neighbours.

This results in the interesting movement of birds/fish/etc. that mimics real life pretty accurately.

In this case, it's schools of fish.

(gif of simulation)

[Complementary rules] and [steering behaviours] allows one to limit the school's speed,
bounding it's position and setting goals etc.

[Boids]: https://en.wikipedia.org/wiki/Boids
[banner]: ./assets/banner.svg
[Craig Reynolds]: https://www.red3d.com/cwr/boids/
[Complementary rules]: https://vergenet.net/~conrad/boids/pseudocode.html
[steering behaviours]: https://gamedevelopment.tutsplus.com/series/understanding-steering-behaviors--gamedev-12732



## Usage

Download the source and then run the main simulation with:

    go run main.go

Or, if you have installed [just], you can simply run it with:

    just run

Running the tests:

    just test

Running all benchmarks for a package:

    just bench boids

Running a specific benchmark inside a package:

    just benchtest=Vectors bench boids

And showing all other shortcut commands:

    just



## FAQ

## Why?

For fun. And I wanted a pretty thing to tinker with once in a while.

### What's the performance like?

Running on my mid-range laptop with an Intel i5-7200U CPU,
I'm able to simulate 10 000 Boids (and 100 goroutine workers) at 60 ± 1 FPS.

Running the benchmark (using 1000 Boids and 10 workers on commit [ce5397c]) I get:

```
# just benchtest=Boids bench boids
go test -bench "Boids" -benchtime 5s -benchmem -cpuprofile ".stats/cpu" -memprofile ".stats/mem" "./boids"
goos: freebsd
goarch: amd64
pkg: github.com/lmas/akvarium/boids
cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
BenchmarkBoids-4   	    2071	   2737401 ns/op	   12814 B/op	       7 allocs/op
PASS
ok  	github.com/lmas/akvarium/boids	6.167s
```


## Issues

- FPS seems to drop when the main window loses focus.

## Roadmap

### Phase One

- Public release and feedback.

### Phase Two

- Add simple entities, moving in the distant background; couple of whales, a giant sunfish.
- Add static environment for back- and foreground; prerendered corals, rocks and cliffs[^1].
- Add shader to simulate some kind of distance blur/dim for background objects.

### Phase Three

- Add camera controls (zoom and pan), so the aquarium can increase in volume.
- Add animation; all entities should render a couple of frames.
- Add more entities; floating jellyfish, shrimp, feeding anemones, swaying kelp stalks.
- Replace the school of clownfish; herring seems more appropriate?
- Add shader to simulate shimmering fish scales.

### Phase Four

- Investigate randomly generating environment, such as corals and rocks[^2].
- Add underwater sounds.
- Add some form of user interaction (and saving state); feeding fish?

## License

GPL, See the [LICENSE] file for details.

[just]: https://github.com/casey/just
[ce5397c]: https://github.com/lmas/akvarium/commit/ce5397cee27cf6f4698a6bcff17b314aaca788b5
[LICENSE]: LICENSE
[Lucas Milner]: https://www.lucasmilner.com/growing-virtual-coral
[space colonization]: http://marcinignac.com/experiments/space-colonization/

[^1]: Should investigate ocean environments properly and decide which one to simulate accurately. Probably tropical.
[^2]: [Lucas Milner] has a great inspirational page. And it seems to be based on a [space colonization] algorithm.
