BINARY=vizid

.PHONY: build test lint run

build:
	go build -o bin/$(BINARY) ./cmd/vizid

run:
	go run ./cmd/vizid gen

test:
	go test ./...

