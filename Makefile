.PHONY: build clean test coverage lint fmt vet sec install help
.DEFAULT_GOAL := help

# Build variables
BINARY_NAME := lfr-tools
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X github.com/scttfrdmn/lfr-tools/cmd.version=$(VERSION) -X github.com/scttfrdmn/lfr-tools/cmd.commit=$(COMMIT) -X github.com/scttfrdmn/lfr-tools/cmd.date=$(DATE)"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

## clean: Remove build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/ dist/ coverage.out

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -race -v ./...

## coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w -local github.com/scttfrdmn/lfr-tools .

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## sec: Run security checks with gosec
sec:
	@echo "Running security checks..."
	@gosec ./...

## mod: Tidy and verify modules
mod:
	@echo "Tidying modules..."
	@go mod tidy
	@go mod verify

## install: Install binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) .

## check: Run all checks (fmt, vet, lint, sec, test)
check: fmt vet lint sec test

## release-dry: Dry run of goreleaser
release-dry:
	@echo "Running goreleaser dry run..."
	@goreleaser release --snapshot --rm-dist --skip-publish

## deps: Install development dependencies
deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/goreleaser/goreleaser@latest

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)