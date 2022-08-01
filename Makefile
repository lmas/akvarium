COVER=.stats/cover
COVER_HTML=.stats/cover.html
CPU=.stats/cpu
MEM=.stats/mem

.PHONY: run
run:
	go run main.go -verbose

.PHONY: test
test:
	go test -coverprofile="${COVER}" $$(go list ./... | grep -v assets | grep -v tools)

.PHONY: cover
cover:
	go tool cover -html="${COVER}" -o="${COVER_HTML}"

.PHONY: benchboids
benchboids:
	go test -bench="Boids" -benchtime=5s -benchmem -cpuprofile="${CPU}" -memprofile="${MEM}" ./boids
	go tool pprof -top "${CPU}"
	rm boids.test

.PHONY: benchvectors
benchvectors:
	go test -bench="Vectors" -benchtime=5s -benchmem -cpuprofile="${CPU}" -memprofile="${MEM}" ./boids
	go tool pprof -top "${CPU}"
	rm boids.test

.PHONY: showcpu
showcpu:
	go tool pprof -http localhost:8000 "${CPU}"

.PHONY: showmem
showmem:
	go tool pprof -http localhost:8000 "${MEM}"

.PHONY: lint
lint:
	golangci-lint run -E gosec -E gocritic --skip-dirs="(assets|tools)" ./...

.PHONY: clean
clean:
	go clean
	rm .stats/*
