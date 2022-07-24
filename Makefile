COVER=.stats/cover
COVER_HTML=.stats/cover.html
CPU=.stats/cpu

.PHONY: run
run:
	go run main.go -debug

.PHONY: test
test:
	go test -coverprofile="${COVER}" ./...

.PHONY: cover
cover:
	go tool cover -html="${COVER}" -o="${COVER_HTML}"

.PHONY: bench
bench:
	go test -bench=. -benchtime=5s -benchmem -cpuprofile="${CPU}" ./simulation

.PHONY: profile
profile: bench
	go tool pprof -top "${CPU}"

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
