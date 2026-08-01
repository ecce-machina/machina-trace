VERSION := $(shell git describe --tags --always --dirty)
COMMIT  := $(shell git rev-parse --short HEAD)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X 'main.Version=$(VERSION)' \
           -X 'main.Commit=$(COMMIT)' \
           -X 'main.BuildDate=$(DATE)'

fmt:
	go fmt ./...

test:
	go test ./...

build:
	go build -ldflags "$(LDFLAGS)" -o bin/machina-agent ./cmd/machina-agent
	go build -ldflags "$(LDFLAGS)" -o bin/machina-trace ./cmd/machina-trace

run:
	go run ./cmd/machina-agent

