VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-s -w -X github.com/jerryagbesi/skipper/cmd.version=$(VERSION)" -o skipper

test:
	go test ./...

run:
	go run .

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

.PHONY: all build run lint fmt
all:
	golangci-lint fmt && go build -o skipper && go run .
