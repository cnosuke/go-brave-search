.PHONY: test cover lint fmt clean build example

# Default target
all: test lint build

# Build the library
build:
	go build ./...

# Run tests
test:
	go test -v ./...

# Run tests with coverage
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run linter
lint:
	go vet ./...
	$(if $(shell which golint), golint ./..., @echo "golint not installed. Run: go install golang.org/x/lint/golint@latest")

# Format code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	rm -f coverage.out
	rm -rf bin/

# Build and run example
example:
	go run examples/simple/main.go

# Install dependencies
deps:
	go mod tidy
	$(if $(shell which golint),,go install golang.org/x/lint/golint@latest)
