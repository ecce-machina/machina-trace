fmt:
	go fmt ./...

test:
	go test ./...

build:
	go build ./cmd/machina-agent
	go build ./cmd/machina-trace

run:
	go run ./cmd/machina-agent

