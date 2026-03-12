.PHONY: build run test docker-build clean help

help:
	@echo "Usage:"
	@echo "  make build         Build the mlc-markitdown binary"
	@echo "  make run           Run the server"
	@echo "  make test          Run tests"
	@echo "  make docker-build  Build Docker image"
	@echo "  make clean         Remove built binaries"

build:
	go build -o bin/mlc-markitdown ./cmd/server/main.go

run:
	go run ./cmd/server/main.go

test:
	go test ./...

docker-build:
	docker build -t mlc-markitdown .

clean:
	rm -rf bin/
