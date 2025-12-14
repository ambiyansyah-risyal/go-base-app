.PHONY: build test run-server run-cli clean

build:
	mkdir -p build
	go build -o build/server ./cmd/server
	go build -o build/cli ./cmd/cli

test:
	go test ./...

run-server:
	go run ./cmd/server

run-cli:
	go run ./cmd/cli

clean:
	rm -rf build
