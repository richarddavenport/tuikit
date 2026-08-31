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

# Pinned to what CI runs, and run through `go run` so it is the SAME version
# whether or not anything is installed. An earlier version of this skipped the
# linter when it was missing, which is how a config CI could not load got
# pushed: a local check that quietly omits what CI enforces is a check that lies.
GOLANGCI ?= go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.0

## lint: what CI enforces
lint:
	gofmt -l . | tee /dev/stderr | (! read)
	go vet ./...
	$(GOLANGCI) run ./...

## check: everything CI runs, before you push
check: test lint
	@echo "all checks passed"

## designsystem: write tuikit's own foundations bundle to ./design-system
designsystem:
	@go run ./cmd/tuikit designsystem -out design-system -tool tuikit

## help: this
help:
	@grep -h '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
