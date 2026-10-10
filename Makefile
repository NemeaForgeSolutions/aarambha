.PHONY: all build run dev test clean css

BINARY_NAME=bin/server

all: build

build:
	@mkdir -p bin
	go build -o $(BINARY_NAME) ./cmd/server

run: build
	./$(BINARY_NAME)

dev:
	go run ./cmd/server

test:
	go test -v ./...

clean:
	rm -rf bin server
	go clean

css:
	npm run build:css

docker-build:
	docker build -t aarambha-agency:latest .

docker-run:
	docker run -p 8080:8080 -e PORT=8080 aarambha-agency:latest

upload-images:
	go run ./cmd/uploader

upload-images-dry:
	go run ./cmd/uploader --dry-run

