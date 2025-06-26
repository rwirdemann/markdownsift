# Makefile for markdownsift

BINARY_NAME=mds
BUILD_DIR=cmd/markdownsift
MAIN_FILE=$(BUILD_DIR)/main.go

# Default target
all: build

# Build the binary
build:
	go build -o $(GOPATH)/bin/$(BINARY_NAME) $(MAIN_FILE)

# Clean built artifacts
clean:
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Run tests
test:
	go test -v

# Run tests with coverage
test-coverage:
	go test -v -cover

.PHONY: all build clean test test-coverage
