# Local development. `make check` is what CI runs and what you run before
# pushing.
#
# tuikit is mostly a library, but `tuikit new` and `tuikit gallery` are things
# you type — so there IS something to install, and there did not used to be.

BIN     ?= $(HOME)/.local/bin/tuikit
VERSION ?= $(shell git describe --tags --dirty 2>/dev/null || echo dev)
LDFLAGS  = -s -w -X main.version=$(VERSION)

.PHONY: install build test lint check designsystem gallery help

## install: build to ~/.local/bin/tuikit (on your PATH)
install:
	@CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o "$(BIN)" ./cmd/tuikit
	@echo "installed $$("$(BIN)" version) -> $(BIN)"

## build: build the tuikit binary to ./bin/tuikit
build:
	@CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tuikit ./cmd/tuikit
	@echo "built $$(./bin/tuikit version) -> bin/tuikit"

## gallery: open every component, running, without installing anything
gallery:
	@go run ./cmd/tuikit gallery

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
