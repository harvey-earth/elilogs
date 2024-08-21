.DEFAULT_GOAL := build

.PHONY: fmt vet build test
fmt:
	go fmt ./...
vet: fmt
	go mod tidy
	go vet ./...
build: manpage
	go build -o bin/elilogs
run: vet
	go run main.go $(arg)
manpage: vet
	go run docs/main.go
prod: vet
	go build -o bin/elilogs
