build:
	go build -o skipper

run:
	go run .

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

.PHONY: all build run lint fmt
all:
	golangci-lint fmt && go build -o skipper && go run .
