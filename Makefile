COVER=.stats/cover
COVER_HTML=.stats/cover.html
CPU=.stats/cpu
MEM=.stats/mem

.PHONY: run
run:
	go run main.go

.PHONY: debug
debug:
	go run main.go -debug=true -pretty=false -init=0

.PHONY: test
test:
	go test -coverprofile="${COVER}" ./...

.PHONY: cover
cover:
	go tool cover -html="${COVER}" -o="${COVER_HTML}"

.PHONY: bench
bench:
	go test -bench=. -benchtime=5s -benchmem -cpuprofile="${CPU}" -memprofile="${MEM}" ./boids
	rm boids.test

.PHONY: profile
profile: bench
	go tool pprof -top "${CPU}"

.PHONY: showcpu
showcpu:
	go tool pprof -http localhost:8000 "${CPU}"

.PHONY: showmem
showmem:
	go tool pprof -http localhost:8000 "${MEM}"

.PHONY: lint
lint:
	golangci-lint run -E gosec -E gocritic ./...

.PHONY: clean
clean:
	go clean
	rm .stats/*
