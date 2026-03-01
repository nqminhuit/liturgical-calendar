.PHONY: all build test compile clean fmt vet run

# Default target
all: test build

# Build the binary
build:
	mkdir -p bin
	go build -ldflags="-s -w" -v -o bin/lc cmd/main.go

# Run tests with detailed report
test:
	go test -v -cover ./...

# Compile code without building binary (check for errors)
compile:
	go build ./...

# Clean up build artifacts
clean:
	rm -rf bin/
	rm -rf liturgical-calendar-*.json

# Format Go code
fmt:
	go fmt ./...

# Run Go vet for static analysis
vet:
	go vet ./...

# Run the program with a specified year
run:
ifndef YEAR
	@echo "Usage: make run YEAR=<year>"
else
	go run cmd/main.go $(YEAR)
endif

# Combined check: format, vet, compile
check: fmt vet compile
