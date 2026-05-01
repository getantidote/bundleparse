export GO111MODULE := "on"

SOURCE_FILES := "./..."
TEST_PATTERN := "."
TEST_OPTIONS := "-v"

# Default target
default: build

# Run all the tests
test:
    go test {{TEST_OPTIONS}} -failfast -race {{SOURCE_FILES}} -run {{TEST_PATTERN}} -timeout=2m

# Run benchmarks from Go
benchmark:
	go test -bench=. -count=1 ./...

# Run all the linters
lint:
    golangci-lint run ./...

# Build a beta version
build:
    go build

# Format all go files
fmt:
    go fmt ./...

# Run all the tests and code checks
ci: build test lint
