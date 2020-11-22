.PHONY: build fmt test

build:
	go build -o bin/mcrawl cmd/mcrawl/main.go

fmt:
	go fmt ./...

test:
	go test -race -v ./...
