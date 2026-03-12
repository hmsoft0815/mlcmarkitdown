.PHONY: build run test docker-build clean help

help:
	@echo "Usage:"
	@echo "  make build         Build the mlc-markitdown binary"
	@echo "  make run           Run the server"
	@echo "  make test          Run tests"
	@echo "  make vet           Run go vet"
	@echo "  make fmt           Run go fmt"
	@echo "  make lint          Run golangci-lint"
	@echo "  make docker-build  Build Docker image"
	@echo "  make clean         Remove built binaries"

build:
	go build -o bin/mlc-markitdown ./cmd/server/main.go

run:
	go run ./cmd/server/main.go

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

docker-build:
	docker build -t mlc-markitdown .

clean:
	rm -rf bin/
