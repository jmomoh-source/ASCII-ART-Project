# ASCII-ART-Project Makefile

# Binary names
BINARY_NAME=ascii-art
VERSION=1.0.0

.PHONY: all build test clean build-all

all: test build

# Compile natively for your OS
build:
	go build -o $(BINARY_NAME) .
	@echo "Built $(BINARY_NAME) successfully."

# Run tests
test:
	go test -v ./...
	@echo "All tests passed successfully."

# Remove built executables
clean:
	rm -f $(BINARY_NAME)* bin/*
	@echo "Cleaned up built files."

# Cross-compile for Major Operating Systems
build-all: clean
	@echo "Building for Linux (x86_64)..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 .
	
	@echo "Building for Windows (x86_64)..."
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 .
	@echo "Cross-compilation complete! Check the bin/ directory."
