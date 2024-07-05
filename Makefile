.PHONY: all test lint

BIN := .bin

all: lint test

test:
	go test ./...

lint: .bin/golangci-lint
	golangci-lint run ./...

${BIN}/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .bin v1.59.1
