.PHONY: run build test

run:
	go run cmd/api/main.go

build:
	go build -o blog-platform cmd/api/main.go

test:
	go test ./... -v