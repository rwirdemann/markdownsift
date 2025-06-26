# Makefile for markdownsift

BINARY_NAME=mds
BUILD_DIR=cmd/markdownsift
MAIN_FILE=$(BUILD_DIR)/main.go

# Default target
all: build

# Build the binary
build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

# Install the binary to $GOPATH/bin
install: build
	mkdir -p $(GOPATH)/bin
	cp $(BINARY_NAME) $(GOPATH)/bin/

# Clean built artifacts
clean:
	rm -f $(BINARY_NAME)

# Run tests
test:
	go test -v

# Run tests with coverage
test-coverage:
	go test -v -cover

# Build and run
run: build
	./$(BINARY_NAME)

# Display help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  install       - Build and install to \$$GOPATH/bin"
	@echo "  clean         - Remove built artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  run           - Build and run the binary"
	@echo "  help          - Show this help message"

.PHONY: all build install clean test test-coverage run help
