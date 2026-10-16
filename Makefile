.PHONY: build test check

build:
	go build -o bin/keyboard-tester .

test:
	go test -race ./...

check:
	go vet ./...
	go test -race ./...
