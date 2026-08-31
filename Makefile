# Local development. tuikit is a library, so there is nothing to install —
# `make check` is what CI runs and what you run before pushing.

VERSION ?= $(shell git describe --tags --dirty 2>/dev/null || echo dev)
LDFLAGS  = -s -w -X main.version=$(VERSION)

.PHONY: build test lint check designsystem help

## build: build the tuikit binary to ./bin/tuikit
build:
	@CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tuikit ./cmd/tuikit
	@echo "built $$(./bin/tuikit version) -> bin/tuikit"

## test: unit tests
test:
	go test ./...

## lint: what CI enforces
lint:
	gofmt -l . | tee /dev/stderr | (! read)
	go vet ./...
	@command -v golangci-lint >/dev/null && golangci-lint run ./... || echo "golangci-lint not installed, skipped"

## check: everything CI runs, before you push
check: test lint
	@echo "all checks passed"

## designsystem: write tuikit's own foundations bundle to ./design-system
designsystem:
	@go run ./cmd/tuikit designsystem -out design-system -tool tuikit

## help: this
help:
	@grep -h '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
