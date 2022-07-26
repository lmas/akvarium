COVER=.stats/cover
COVER_HTML=.stats/cover.html
CPU=.stats/cpu
MEM=.stats/mem

.PHONY: run
run:
	go run main.go

.PHONY: debug
debug:
	go run main.go -debug=true -effects=false

.PHONY: test
test:
	go test -coverprofile="${COVER}" ./...

.PHONY: cover
cover:
	go tool cover -html="${COVER}" -o="${COVER_HTML}"

.PHONY: bench
bench:
	go test -bench=. -benchtime=5s -benchmem -cpuprofile="${CPU}" -memprofile="${MEM}" ./simulation
	rm simulation.test

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
lint: govet gosec

.PHONY: govet
govet:
	go vet ./...

.PHONY: gosec
gosec:
	gosec -quiet -fmt=golint ./...

.PHONY: clean
clean:
	go clean
	rm .stats/*
