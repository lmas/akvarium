STATS := ".stats"
IGNORE := "(assets|tools)"

# Show available recipes by default
default:
    @just --list

# Run the main simulation
run:
    go run main.go -verbose

# Run linter on source code
lint:
    golangci-lint run -E gosec -E gocritic --skip-dirs "{{IGNORE}}" ./...

# Run all tests
test:
   go test -coverprofile "{{STATS}}/cover" $(go list ./... | grep -v -E "{{IGNORE}}")

# Produce coverage web report
cover:
    go tool cover -html "{{STATS}}/cover" -o "{{STATS}}/cover.html"

benchtest := "." # limit what benchmarks to run, eg: just benchtest=... bench ...

# Run specific benchmark
bench target:
    go test -bench "{{benchtest}}" -benchtime 5s -benchmem \
    -cpuprofile "{{STATS}}/cpu" -memprofile "{{STATS}}/mem" "./{{target}}"
    go tool pprof -top "{{STATS}}/cpu"
    rm "{{target}}.test"

# Run pprof web explorer
pprof target:
    go tool pprof -http localhost:8000 "{{STATS}}/{{target}}"

# Cleanup artefacts
clean:
    go clean
    rm {{STATS}}/*
